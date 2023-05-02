import { cloneDeep, find, first as _first, isObject, isString, map as _map } from 'lodash';
import { generate, lastValueFrom, Observable, of } from 'rxjs';
import { catchError, first, map, mergeMap, skipWhile, throwIfEmpty, tap } from 'rxjs/operators';
import { SemVer } from 'semver';

import {
  DataFrame,
  DataLink,
  DataQueryRequest,
  DataQueryResponse,
  DataSourceInstanceSettings,
  DataSourceWithLogsContextSupport,
  DataSourceWithQueryImportSupport,
  DataSourceWithSupplementaryQueriesSupport,
  DateTime,
  dateTime,
  getDefaultTimeRange,
  AbstractQuery,
  LogLevel,
  LogRowModel,
  MetricFindValue,
  ScopedVars,
  TimeRange,
  QueryFixAction,
  CoreApp,
  SupplementaryQueryType,
  DataQueryError,
  rangeUtil,
  LogRowContextQueryDirection,
  LogRowContextOptions,
  SupplementaryQueryOptions,
} from '@grafana/data';
import { DataSourceWithBackend, getDataSourceSrv, config } from '@grafana/runtime';
import { queryLogsVolume } from 'app/core/logsModel';
import { getTimeSrv, TimeSrv } from 'app/features/dashboard/services/TimeSrv';
import { getTemplateSrv, TemplateSrv } from 'app/features/templating/template_srv';

import { getLogLevelFromKey } from '../../../features/logs/utils';

import { IndexPattern, intervalMap } from './IndexPattern';
import LanguageProvider from './LanguageProvider';
import { LegacyQueryRunner } from './LegacyQueryRunner';
import { ElasticQueryBuilder } from './QueryBuilder';
import { ElasticsearchAnnotationsQueryEditor } from './components/QueryEditor/AnnotationQueryEditor';
import { isBucketAggregationWithField } from './components/QueryEditor/BucketAggregationsEditor/aggregations';
import { bucketAggregationConfig } from './components/QueryEditor/BucketAggregationsEditor/utils';
import {
  isMetricAggregationWithField,
  isPipelineAggregationWithMultipleBucketPaths,
} from './components/QueryEditor/MetricAggregationsEditor/aggregations';
import { metricAggregationConfig } from './components/QueryEditor/MetricAggregationsEditor/utils';
import { trackQuery } from './tracking';
import {
  Logs,
  BucketAggregation,
  DataLinkConfig,
  ElasticsearchOptions,
  ElasticsearchQuery,
  TermsQuery,
  Interval,
} from './types';
import { getScriptValue, isSupportedVersion, unsupportedVersionMessage } from './utils';

export const REF_ID_STARTER_LOG_VOLUME = 'log-volume-';
// Those are metadata fields as defined in https://www.elastic.co/guide/en/elasticsearch/reference/current/mapping-fields.html#_identity_metadata_fields.
// custom fields can start with underscores, therefore is not safe to exclude anything that starts with one.
const ELASTIC_META_FIELDS = [
  '_index',
  '_type',
  '_id',
  '_source',
  '_size',
  '_field_names',
  '_ignored',
  '_routing',
  '_meta',
];

export class ElasticDatasource
  extends DataSourceWithBackend<ElasticsearchQuery, ElasticsearchOptions>
  implements
    DataSourceWithLogsContextSupport,
    DataSourceWithQueryImportSupport<ElasticsearchQuery>,
    DataSourceWithSupplementaryQueriesSupport<ElasticsearchQuery>
{
  basicAuth?: string;
  withCredentials?: boolean;
  url: string;
  name: string;
  index: string;
  timeField: string;
  xpack: boolean;
  interval: string;
  maxConcurrentShardRequests?: number;
  queryBuilder: ElasticQueryBuilder;
  indexPattern: IndexPattern;
  intervalPattern?: Interval;
  logMessageField?: string;
  logLevelField?: string;
  dataLinks: DataLinkConfig[];
  languageProvider: LanguageProvider;
  includeFrozen: boolean;
  isProxyAccess: boolean;
  timeSrv: TimeSrv;
  databaseVersion: SemVer | null;
  legacyQueryRunner: LegacyQueryRunner;

  constructor(
    instanceSettings: DataSourceInstanceSettings<ElasticsearchOptions>,
    private readonly templateSrv: TemplateSrv = getTemplateSrv()
  ) {
    super(instanceSettings);
    this.basicAuth = instanceSettings.basicAuth;
    this.withCredentials = instanceSettings.withCredentials;
    this.url = instanceSettings.url!;
    this.name = instanceSettings.name;
    this.isProxyAccess = instanceSettings.access === 'proxy';
    const settingsData = instanceSettings.jsonData || ({} as ElasticsearchOptions);

    this.index = settingsData.index ?? instanceSettings.database ?? '';
    this.timeField = settingsData.timeField;
    this.xpack = Boolean(settingsData.xpack);
    this.indexPattern = new IndexPattern(this.index, settingsData.interval);
    this.intervalPattern = settingsData.interval;
    this.interval = settingsData.timeInterval;
    this.maxConcurrentShardRequests = settingsData.maxConcurrentShardRequests;
    this.queryBuilder = new ElasticQueryBuilder({
      timeField: this.timeField,
    });
    this.logMessageField = settingsData.logMessageField || '';
    this.logLevelField = settingsData.logLevelField || '';
    this.dataLinks = settingsData.dataLinks || [];
    this.includeFrozen = settingsData.includeFrozen ?? false;
    this.databaseVersion = null;
    this.annotations = {
      QueryEditor: ElasticsearchAnnotationsQueryEditor,
    };

    if (this.logMessageField === '') {
      this.logMessageField = undefined;
    }

    if (this.logLevelField === '') {
      this.logLevelField = undefined;
    }
    this.languageProvider = new LanguageProvider(this);
    this.timeSrv = getTimeSrv();
    this.legacyQueryRunner = new LegacyQueryRunner(this, this.templateSrv);
  }

  async importFromAbstractQueries(abstractQueries: AbstractQuery[]): Promise<ElasticsearchQuery[]> {
    return abstractQueries.map((abstractQuery) => this.languageProvider.importFromAbstractQuery(abstractQuery));
  }

  /**
   * Sends a GET request to the specified url on the newest matching and available index.
   *
   * When multiple indices span the provided time range, the request is sent starting from the newest index,
   * and then going backwards until an index is found.
   *
   * @param url the url to query the index on, for example `/_mapping`.
   */

  private requestAllIndices(url: string, range = getDefaultTimeRange()): Observable<any> {
    let indexList = this.indexPattern.getIndexList(range.from, range.to);
    if (!Array.isArray(indexList)) {
      indexList = [this.indexPattern.getIndexForToday()];
    }

    const indexUrlList = indexList.map((index) => index + url);

    const maxTraversals = 7; // do not go beyond one week (for a daily pattern)
    const listLen = indexUrlList.length;

    return generate({
      initialState: 0,
      condition: (i) => i < Math.min(listLen, maxTraversals),
      iterate: (i) => i + 1,
    }).pipe(
      mergeMap((index) => {
        // catch all errors and emit an object with an err property to simplify checks later in the pipeline
        return this.legacyQueryRunner
          .request('GET', indexUrlList[listLen - index - 1])
          .pipe(catchError((err) => of({ err })));
      }),
      skipWhile((resp) => resp?.err?.status === 404), // skip all requests that fail because missing Elastic index
      throwIfEmpty(() => 'Could not find an available index for this time range.'), // when i === Math.min(listLen, maxTraversals) generate will complete but without emitting any values which means we didn't find a valid index
      first(), // take the first value that isn't skipped
      map((resp) => {
        if (resp.err) {
          throw resp.err; // if there is some other error except 404 then we must throw it
        }

        return resp;
      })
    );
  }

  annotationQuery(options: any): Promise<any> {
    return this.legacyQueryRunner.annotationQuery(options);
  }

  interpolateLuceneQuery(queryString: string, scopedVars?: ScopedVars) {
    return this.templateSrv.replace(queryString, scopedVars, 'lucene');
  }

  interpolateVariablesInQueries(queries: ElasticsearchQuery[], scopedVars: ScopedVars | {}): ElasticsearchQuery[] {
    return queries.map((q) => this.applyTemplateVariables(q, scopedVars));
  }

  async testDatasource() {
    // we explicitly ask for uncached, "fresh" data here
    const dbVersion = await this.getDatabaseVersion(false);
    // if we are not able to determine the elastic-version, we assume it is a good version.
    const isSupported = dbVersion != null ? isSupportedVersion(dbVersion) : true;
    const versionMessage = isSupported ? '' : `WARNING: ${unsupportedVersionMessage} `;
    // validate that the index exist and has date field
    return lastValueFrom(
      this.getFields(['date']).pipe(
        mergeMap((dateFields) => {
          const timeField: any = find(dateFields, { text: this.timeField });
          if (!timeField) {
            return of({
              status: 'error',
              message: 'No date field named ' + this.timeField + ' found',
            });
          }
          return of({ status: 'success', message: `${versionMessage}Index OK. Time field name OK` });
        }),
        catchError((err) => {
          console.error(err);
          if (err.message) {
            return of({ status: 'error', message: err.message });
          } else {
            return of({ status: 'error', message: err.status });
          }
        })
      )
    );
  }

  getQueryHeader(searchType: any, timeFrom?: DateTime, timeTo?: DateTime): string {
    const queryHeader: any = {
      search_type: searchType,
      ignore_unavailable: true,
      index: this.indexPattern.getIndexList(timeFrom, timeTo),
    };

    return JSON.stringify(queryHeader);
  }

  getQueryDisplayText(query: ElasticsearchQuery) {
    // TODO: This might be refactored a bit.
    const metricAggs = query.metrics;
    const bucketAggs = query.bucketAggs;
    let text = '';

    if (query.query) {
      text += 'Query: ' + query.query + ', ';
    }

    text += 'Metrics: ';

    text += metricAggs?.reduce((acc, metric) => {
      const metricConfig = metricAggregationConfig[metric.type];

      let text = metricConfig.label + '(';

      if (isMetricAggregationWithField(metric)) {
        text += metric.field;
      }
      if (isPipelineAggregationWithMultipleBucketPaths(metric)) {
        text += getScriptValue(metric).replace(new RegExp('params.', 'g'), '');
      }
      text += '), ';

      return `${acc} ${text}`;
    }, '');

    text += bucketAggs?.reduce((acc, bucketAgg, index) => {
      const bucketConfig = bucketAggregationConfig[bucketAgg.type];

      let text = '';
      if (index === 0) {
        text += ' Group by: ';
      }

      text += bucketConfig.label + '(';
      if (isBucketAggregationWithField(bucketAgg)) {
        text += bucketAgg.field;
      }

      return `${acc} ${text}), `;
    }, '');

    if (query.alias) {
      text += 'Alias: ' + query.alias;
    }

    return text;
  }

  showContextToggle(): boolean {
    return true;
  }

  getLogRowContext = async (row: LogRowModel, options?: LogRowContextOptions): Promise<{ data: DataFrame[] }> => {
    const { enableElasticsearchBackendQuerying } = config.featureToggles;
    if (enableElasticsearchBackendQuerying) {
      const contextRequest = this.makeLogContextDataRequest(row, options);

      return lastValueFrom(
        this.query(contextRequest).pipe(
          catchError((err) => {
            const error: DataQueryError = {
              message: 'Error during context query. Please check JS console logs.',
              status: err.status,
              statusText: err.statusText,
            };
            throw error;
          })
        )
      );
    } else {
      return this.legacyQueryRunner.logContextQuery(row, options);
    }
  };

  getDataProvider(
    type: SupplementaryQueryType,
    request: DataQueryRequest<ElasticsearchQuery>
  ): Observable<DataQueryResponse> | undefined {
    if (!this.getSupportedSupplementaryQueryTypes().includes(type)) {
      return undefined;
    }
    switch (type) {
      case SupplementaryQueryType.LogsVolume:
        return this.getLogsVolumeDataProvider(request);
      default:
        return undefined;
    }
  }

  getSupportedSupplementaryQueryTypes(): SupplementaryQueryType[] {
    return [SupplementaryQueryType.LogsVolume];
  }

  getSupplementaryQuery(options: SupplementaryQueryOptions, query: ElasticsearchQuery): ElasticsearchQuery | undefined {
    if (!this.getSupportedSupplementaryQueryTypes().includes(options.type)) {
      return undefined;
    }

    let isQuerySuitable = false;

    switch (options.type) {
      case SupplementaryQueryType.LogsVolume:
        // it has to be a logs-producing range-query
        isQuerySuitable = !!(query.metrics?.length === 1 && query.metrics[0].type === 'logs');
        if (!isQuerySuitable) {
          return undefined;
        }
        const bucketAggs: BucketAggregation[] = [];
        const timeField = this.timeField ?? '@timestamp';

        if (this.logLevelField) {
          bucketAggs.push({
            id: '2',
            type: 'terms',
            settings: {
              min_doc_count: '0',
              size: '0',
              order: 'desc',
              orderBy: '_count',
              missing: LogLevel.unknown,
            },
            field: this.logLevelField,
          });
        }
        bucketAggs.push({
          id: '3',
          type: 'date_histogram',
          settings: {
            interval: 'auto',
            min_doc_count: '0',
            trimEdges: '0',
          },
          field: timeField,
        });

        return {
          refId: `${REF_ID_STARTER_LOG_VOLUME}${query.refId}`,
          query: query.query,
          metrics: [{ type: 'count', id: '1' }],
          timeField,
          bucketAggs,
        };

      default:
        return undefined;
    }
  }

  getLogsVolumeDataProvider(request: DataQueryRequest<ElasticsearchQuery>): Observable<DataQueryResponse> | undefined {
    const logsVolumeRequest = cloneDeep(request);
    const targets = logsVolumeRequest.targets
      .map((target) => this.getSupplementaryQuery({ type: SupplementaryQueryType.LogsVolume }, target))
      .filter((query): query is ElasticsearchQuery => !!query);

    if (!targets.length) {
      return undefined;
    }

    return queryLogsVolume(
      this,
      { ...logsVolumeRequest, targets },
      {
        range: request.range,
        targets: request.targets,
        extractLevel: (dataFrame) => getLogLevelFromKey(dataFrame.name || ''),
      }
    );
  }

  query(request: DataQueryRequest<ElasticsearchQuery>): Observable<DataQueryResponse> {
    const { enableElasticsearchBackendQuerying } = config.featureToggles;
    if (enableElasticsearchBackendQuerying) {
      const start = new Date();
      return super.query(request).pipe(tap((response) => trackQuery(response, request, start)));
    }
    return this.legacyQueryRunner.query(request);
  }

  isMetadataField(fieldName: string) {
    return ELASTIC_META_FIELDS.includes(fieldName);
  }

  // TODO: instead of being a string, this could be a custom type representing all the elastic types
  // FIXME: This doesn't seem to return actual MetricFindValues, we should either change the return type
  // or fix the implementation.
  getFields(type?: string[], range?: TimeRange): Observable<MetricFindValue[]> {
    const typeMap: Record<string, string> = {
      float: 'number',
      double: 'number',
      integer: 'number',
      long: 'number',
      date: 'date',
      date_nanos: 'date',
      string: 'string',
      text: 'string',
      scaled_float: 'number',
      nested: 'nested',
      histogram: 'number',
    };
    return this.requestAllIndices('/_mapping', range).pipe(
      map((result) => {
        const shouldAddField = (obj: any, key: string) => {
          if (this.isMetadataField(key)) {
            return false;
          }

          if (!type || type.length === 0) {
            return true;
          }

          // equal query type filter, or via typemap translation
          return type.includes(obj.type) || type.includes(typeMap[obj.type]);
        };

        // Store subfield names: [system, process, cpu, total] -> system.process.cpu.total
        const fieldNameParts: any = [];
        const fields: any = {};

        function getFieldsRecursively(obj: any) {
          for (const key in obj) {
            const subObj = obj[key];

            // Check mapping field for nested fields
            if (isObject(subObj.properties)) {
              fieldNameParts.push(key);
              getFieldsRecursively(subObj.properties);
            }

            if (isObject(subObj.fields)) {
              fieldNameParts.push(key);
              getFieldsRecursively(subObj.fields);
            }

            if (isString(subObj.type)) {
              const fieldName = fieldNameParts.concat(key).join('.');

              // Hide meta-fields and check field type
              if (shouldAddField(subObj, key)) {
                fields[fieldName] = {
                  text: fieldName,
                  type: subObj.type,
                };
              }
            }
          }
          fieldNameParts.pop();
        }

        for (const indexName in result) {
          const index = result[indexName];
          if (index && index.mappings) {
            const mappings = index.mappings;

            const properties = mappings.properties;
            getFieldsRecursively(properties);
          }
        }

        // transform to array
        return _map(fields, (value) => {
          return value;
        });
      })
    );
  }

  getTerms(queryDef: TermsQuery, range = getDefaultTimeRange()): Observable<MetricFindValue[]> {
    const searchType = 'query_then_fetch';
    const header = this.getQueryHeader(searchType, range.from, range.to);
    let esQuery = JSON.stringify(this.queryBuilder.getTermsQuery(queryDef));

    esQuery = esQuery.replace(/\$timeFrom/g, range.from.valueOf().toString());
    esQuery = esQuery.replace(/\$timeTo/g, range.to.valueOf().toString());
    esQuery = header + '\n' + esQuery + '\n';

    const url = this.getMultiSearchUrl();

    return this.legacyQueryRunner.request('POST', url, esQuery).pipe(
      map((res) => {
        if (!res.responses[0].aggregations) {
          return [];
        }

        const buckets = res.responses[0].aggregations['1'].buckets;
        return _map(buckets, (bucket) => {
          return {
            text: bucket.key_as_string || bucket.key,
            value: bucket.key,
          };
        });
      })
    );
  }

  getMultiSearchUrl() {
    const searchParams = new URLSearchParams();

    if (this.maxConcurrentShardRequests) {
      searchParams.append('max_concurrent_shard_requests', `${this.maxConcurrentShardRequests}`);
    }

    if (this.xpack && this.includeFrozen) {
      searchParams.append('ignore_throttled', 'false');
    }

    return ('_msearch?' + searchParams.toString()).replace(/\?$/, '');
  }

  metricFindQuery(query: string, options?: any): Promise<MetricFindValue[]> {
    const range = options?.range;
    const parsedQuery = JSON.parse(query);
    if (query) {
      if (parsedQuery.find === 'fields') {
        parsedQuery.type = this.interpolateLuceneQuery(parsedQuery.type);
        return lastValueFrom(this.getFields(parsedQuery.type, range));
      }

      if (parsedQuery.find === 'terms') {
        parsedQuery.field = this.interpolateLuceneQuery(parsedQuery.field);
        parsedQuery.query = this.interpolateLuceneQuery(parsedQuery.query);
        return lastValueFrom(this.getTerms(parsedQuery, range));
      }
    }

    return Promise.resolve([]);
  }

  getTagKeys() {
    return lastValueFrom(this.getFields());
  }

  getTagValues(options: any) {
    const range = this.timeSrv.timeRange();
    return lastValueFrom(this.getTerms({ field: options.key }, range));
  }

  targetContainsTemplate(target: any) {
    if (this.templateSrv.containsTemplate(target.query) || this.templateSrv.containsTemplate(target.alias)) {
      return true;
    }

    for (const bucketAgg of target.bucketAggs) {
      if (this.templateSrv.containsTemplate(bucketAgg.field) || this.objectContainsTemplate(bucketAgg.settings)) {
        return true;
      }
    }

    for (const metric of target.metrics) {
      if (
        this.templateSrv.containsTemplate(metric.field) ||
        this.objectContainsTemplate(metric.settings) ||
        this.objectContainsTemplate(metric.meta)
      ) {
        return true;
      }
    }

    return false;
  }

  private isPrimitive(obj: any) {
    if (obj === null || obj === undefined) {
      return true;
    }
    if (['string', 'number', 'boolean'].some((type) => type === typeof true)) {
      return true;
    }

    return false;
  }

  private objectContainsTemplate(obj: any) {
    if (!obj) {
      return false;
    }

    for (const key of Object.keys(obj)) {
      if (this.isPrimitive(obj[key])) {
        if (this.templateSrv.containsTemplate(obj[key])) {
          return true;
        }
      } else if (Array.isArray(obj[key])) {
        for (const item of obj[key]) {
          if (this.objectContainsTemplate(item)) {
            return true;
          }
        }
      } else {
        if (this.objectContainsTemplate(obj[key])) {
          return true;
        }
      }
    }

    return false;
  }

  modifyQuery(query: ElasticsearchQuery, action: QueryFixAction): ElasticsearchQuery {
    if (!action.options) {
      return query;
    }

    let expression = query.query ?? '';
    switch (action.type) {
      case 'ADD_FILTER': {
        if (expression.length > 0) {
          expression += ' AND ';
        }
        expression += `${action.options.key}:"${action.options.value}"`;
        break;
      }
      case 'ADD_FILTER_OUT': {
        if (expression.length > 0) {
          expression += ' AND ';
        }
        expression += `-${action.options.key}:"${action.options.value}"`;
        break;
      }
    }
    return { ...query, query: expression };
  }

  addAdHocFilters(query: string) {
    const adhocFilters = this.templateSrv.getAdhocFilters(this.name);
    if (adhocFilters.length === 0) {
      return query;
    }
    const esFilters = adhocFilters.map((filter) => {
      const { key, operator, value } = filter;
      if (!key || !value) {
        return;
      }
      switch (operator) {
        case '=':
          return `${key}:"${value}"`;
        case '!=':
          return `-${key}:"${value}"`;
        case '=~':
          return `${key}:/${value}/`;
        case '!~':
          return `-${key}:/${value}/`;
        case '>':
          return `${key}:>${value}`;
        case '<':
          return `${key}:<${value}`;
      }
      return;
    });

    const finalQuery = [query, ...esFilters].filter((f) => f).join(' AND ');
    return finalQuery;
  }

  // Used when running queries through backend
  applyTemplateVariables(query: ElasticsearchQuery, scopedVars: ScopedVars): ElasticsearchQuery {
    // We need a separate interpolation format for lucene queries, therefore we first interpolate any
    // lucene query string and then everything else
    const interpolateBucketAgg = (bucketAgg: BucketAggregation): BucketAggregation => {
      if (bucketAgg.type === 'filters') {
        return {
          ...bucketAgg,
          settings: {
            ...bucketAgg.settings,
            filters: bucketAgg.settings?.filters?.map((filter) => ({
              ...filter,
              query: this.interpolateLuceneQuery(filter.query, scopedVars) || '*',
            })),
          },
        };
      }

      return bucketAgg;
    };

    const expandedQuery = {
      ...query,
      datasource: this.getRef(),
      query: this.addAdHocFilters(this.interpolateLuceneQuery(query.query || '', scopedVars)),
      bucketAggs: query.bucketAggs?.map(interpolateBucketAgg),
    };

    const finalQuery = JSON.parse(this.templateSrv.replace(JSON.stringify(expandedQuery), scopedVars));
    return finalQuery;
  }

  private getDatabaseVersionUncached(): Promise<SemVer | null> {
    // we want this function to never fail
    return lastValueFrom(this.legacyQueryRunner.request('GET', '/')).then(
      (data) => {
        const versionNumber = data?.version?.number;
        if (typeof versionNumber !== 'string') {
          return null;
        }
        try {
          return new SemVer(versionNumber);
        } catch (error) {
          console.error(error);
          return null;
        }
      },
      (error) => {
        console.error(error);
        return null;
      }
    );
  }

  async getDatabaseVersion(useCachedData = true): Promise<SemVer | null> {
    if (useCachedData) {
      const cached = this.databaseVersion;
      if (cached != null) {
        return cached;
      }
    }

    const freshDatabaseVersion = await this.getDatabaseVersionUncached();
    this.databaseVersion = freshDatabaseVersion;
    return freshDatabaseVersion;
  }

  private makeLogContextDataRequest = (row: LogRowModel, options?: LogRowContextOptions) => {
    const direction = options?.direction || LogRowContextQueryDirection.Backward;
    const logQuery: Logs = {
      type: 'logs',
      id: '1',
      settings: {
        limit: options?.limit ? options?.limit.toString() : '10',
        // Sorting of results in the context query
        sortDirection: direction === LogRowContextQueryDirection.Backward ? 'desc' : 'asc',
        // Used to get the next log lines before/after the current log line using sort field of selected log line
        searchAfter: row.dataFrame.fields.find((f) => f.name === 'sort')?.values[row.rowIndex] ?? [row.timeEpochMs],
      },
    };

    const query: ElasticsearchQuery = {
      refId: `log-context-${row.dataFrame.refId}-${direction}`,
      metrics: [logQuery],
      query: '',
    };

    const timeRange = createContextTimeRange(row.timeEpochMs, direction, this.intervalPattern);
    const range = {
      from: timeRange.from,
      to: timeRange.to,
      raw: timeRange,
    };

    const interval = rangeUtil.calculateInterval(range, 1);

    const contextRequest: DataQueryRequest<ElasticsearchQuery> = {
      requestId: `log-context-request-${row.dataFrame.refId}-${options?.direction}`,
      targets: [query],
      interval: interval.interval,
      intervalMs: interval.intervalMs,
      range,
      scopedVars: {},
      timezone: 'UTC',
      app: CoreApp.Explore,
      startTime: Date.now(),
      hideFromInspector: true,
    };
    return contextRequest;
  };
}

/**
 * Modifies dataframe and adds dataLinks from the config.
 * Exported for tests.
 */
export function enhanceDataFrame(dataFrame: DataFrame, dataLinks: DataLinkConfig[], limit?: number) {
  if (limit) {
    dataFrame.meta = {
      ...dataFrame.meta,
      limit,
    };
  }

  if (!dataLinks.length) {
    return;
  }

  for (const field of dataFrame.fields) {
    const linksToApply = dataLinks.filter((dataLink) => new RegExp(dataLink.field).test(field.name));

    if (linksToApply.length === 0) {
      continue;
    }

    field.config = field.config || {};
    field.config.links = [...(field.config.links || [], linksToApply.map(generateDataLink))];
  }
}

function generateDataLink(linkConfig: DataLinkConfig): DataLink {
  const dataSourceSrv = getDataSourceSrv();

  if (linkConfig.datasourceUid) {
    const dsSettings = dataSourceSrv.getInstanceSettings(linkConfig.datasourceUid);

    return {
      title: linkConfig.urlDisplayLabel || '',
      url: '',
      internal: {
        query: { query: linkConfig.url },
        datasourceUid: linkConfig.datasourceUid,
        datasourceName: dsSettings?.name ?? 'Data source not found',
      },
    };
  } else {
    return {
      title: linkConfig.urlDisplayLabel || '',
      url: linkConfig.url,
    };
  }
}

function createContextTimeRange(rowTimeEpochMs: number, direction: string, intervalPattern: Interval | undefined) {
  const offset = 7;
  // For log context, we want to request data from 7 subsequent/previous indices
  if (intervalPattern) {
    const intervalInfo = intervalMap[intervalPattern];
    if (direction === LogRowContextQueryDirection.Forward) {
      return {
        from: dateTime(rowTimeEpochMs).utc(),
        to: dateTime(rowTimeEpochMs).add(offset, intervalInfo.amount).utc().startOf(intervalInfo.startOf),
      };
    } else {
      return {
        from: dateTime(rowTimeEpochMs).subtract(offset, intervalInfo.amount).utc().startOf(intervalInfo.startOf),
        to: dateTime(rowTimeEpochMs).utc(),
      };
    }
    // If we don't have an interval pattern, we can't do this, so we just request data from 7h before/after
  } else {
    if (direction === LogRowContextQueryDirection.Forward) {
      return {
        from: dateTime(rowTimeEpochMs).utc(),
        to: dateTime(rowTimeEpochMs).add(offset, 'hours').utc(),
      };
    } else {
      return {
        from: dateTime(rowTimeEpochMs).subtract(offset, 'hours').utc(),
        to: dateTime(rowTimeEpochMs).utc(),
      };
    }
  }
}