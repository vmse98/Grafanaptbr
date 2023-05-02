import { BaseQueryFn, createApi } from '@reduxjs/toolkit/query/react';
import { lastValueFrom } from 'rxjs';

import { isTruthy } from '@grafana/data';
import { BackendSrvRequest, getBackendSrv } from '@grafana/runtime';
import { DeleteDashboardResponse } from 'app/features/manage-dashboards/types';
import { DashboardDTO, FolderDTO } from 'app/types';

import { DashboardTreeSelection } from '../types';

interface RequestOptions extends BackendSrvRequest {
  manageError?: (err: unknown) => { error: unknown };
  showErrorAlert?: boolean;
}

function createBackendSrvBaseQuery({ baseURL }: { baseURL: string }): BaseQueryFn<RequestOptions> {
  async function backendSrvBaseQuery(requestOptions: RequestOptions) {
    try {
      const { data: responseData, ...meta } = await lastValueFrom(
        getBackendSrv().fetch({
          ...requestOptions,
          url: baseURL + requestOptions.url,
          showErrorAlert: requestOptions.showErrorAlert,
        })
      );
      return { data: responseData, meta };
    } catch (error) {
      return requestOptions.manageError ? requestOptions.manageError(error) : { error };
    }
  }

  return backendSrvBaseQuery;
}

export const browseDashboardsAPI = createApi({
  reducerPath: 'browseDashboardsAPI',
  baseQuery: createBackendSrvBaseQuery({ baseURL: '/api' }),
  endpoints: (builder) => ({
    deleteDashboard: builder.mutation<DeleteDashboardResponse, string>({
      query: (dashboardUID) => ({
        url: `/dashboards/uid/${dashboardUID}`,
        method: 'DELETE',
      }),
    }),
    deleteFolder: builder.mutation<void, string>({
      query: (folderUID) => ({
        url: `/folders/${folderUID}`,
        method: 'DELETE',
        params: {
          forceDeleteRules: true,
        },
      }),
    }),
    // TODO we can define this return type properly
    moveDashboard: builder.mutation<
      unknown,
      {
        dashboardUID: string;
        destinationUID: string;
      }
    >({
      queryFn: async ({ dashboardUID, destinationUID }, _api, _extraOptions, baseQuery) => {
        const fullDash: DashboardDTO = await getBackendSrv().get(`/api/dashboards/uid/${dashboardUID}`);

        const options = {
          dashboard: fullDash.dashboard,
          folderUid: destinationUID,
          overwrite: false,
        };

        return baseQuery({
          url: '/dashboards/db',
          method: 'POST',
          data: {
            message: '',
            ...options,
          },
        });
      },
    }),
    // TODO this doesn't return void, find where the correct type is
    moveFolder: builder.mutation<
      void,
      {
        folderUID: string;
        destinationUID: string;
      }
    >({
      query: ({ folderUID, destinationUID }) => ({
        url: `/folders/${folderUID}/move`,
        method: 'POST',
        data: { parentUid: destinationUID },
      }),
    }),
    getFolder: builder.query<FolderDTO, string>({
      query: (folderUID) => ({ url: `/folders/${folderUID}` }),
    }),
    getAffectedItems: builder.query<
      // TODO move to folder types file once structure is finalised
      {
        folder: number;
        dashboard: number;
        libraryPanel: number;
        alertRule: number;
      },
      DashboardTreeSelection
    >({
      queryFn: async (selectedItems) => {
        const folderUIDs = Object.keys(selectedItems.folder).filter((uid) => selectedItems.folder[uid]);
        // Mock descendant count
        // TODO convert to real implementation
        const mockDescendantCount = {
          folder: 1,
          dashboard: 1,
          libraryPanel: 1,
          alertRule: 1,
        };
        const promises = folderUIDs.map((id) => {
          return new Promise<typeof mockDescendantCount>((resolve, reject) => {
            // Artificial delay to simulate network request
            setTimeout(() => {
              resolve(mockDescendantCount);
              // reject(new Error('Uh oh!'));
            }, 1000);
          });
        });

        const results = await Promise.all(promises);
        const aggregatedResults = results.reduce(
          (acc, val) => ({
            folder: acc.folder + val.folder,
            dashboard: acc.dashboard + val.dashboard,
            libraryPanel: acc.libraryPanel + val.libraryPanel,
            alertRule: acc.alertRule + val.alertRule,
          }),
          {
            folder: 0,
            dashboard: 0,
            libraryPanel: 0,
            alertRule: 0,
          }
        );

        // Add in the top level selected items
        aggregatedResults.folder += Object.values(selectedItems.folder).filter(isTruthy).length;
        aggregatedResults.dashboard += Object.values(selectedItems.dashboard).filter(isTruthy).length;
        return { data: aggregatedResults };
      },
    }),
  }),
});

export const {
  useDeleteDashboardMutation,
  useDeleteFolderMutation,
  useGetAffectedItemsQuery,
  useGetFolderQuery,
  useLazyGetFolderQuery,
  useMoveDashboardMutation,
  useMoveFolderMutation,
} = browseDashboardsAPI;
export { skipToken } from '@reduxjs/toolkit/query/react';
