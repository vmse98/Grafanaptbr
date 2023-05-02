import { css } from '@emotion/css';
import { isEmpty } from 'lodash';
import React, { useEffect } from 'react';
import { connect, ConnectedProps } from 'react-redux';

import { GrafanaTheme2 } from '@grafana/data';
import { useStyles2 } from '@grafana/ui';
import { Page } from 'app/core/components/Page/Page';
import { StoreState } from 'app/types';

import ConfigureAuthCTA from './components/ConfigureAuthCTA';
import { ProviderCard } from './components/ProviderCard';
import { loadSettings } from './state/actions';
import { getProviderUrl } from './utils';

import { getRegisteredAuthProviders } from '.';

interface OwnProps {}

export type Props = OwnProps & ConnectedProps<typeof connector>;

function mapStateToProps(state: StoreState) {
  const { isLoading, providerStatuses } = state.authConfig;
  return {
    isLoading,
    providerStatuses,
  };
}

const mapDispatchToProps = {
  loadSettings,
};

const connector = connect(mapStateToProps, mapDispatchToProps);

export const AuthConfigPageUnconnected = ({ providerStatuses, isLoading, loadSettings }: Props): JSX.Element => {
  const styles = useStyles2(getStyles);

  useEffect(() => {
    loadSettings();
  }, [loadSettings]);

  const authProviders = getRegisteredAuthProviders();
  const enabledProviders = authProviders.filter((p) => providerStatuses[p.id]?.enabled);
  const configuresProviders = authProviders.filter(
    (p) => providerStatuses[p.id]?.configured && !providerStatuses[p.id]?.enabled
  );
  const availableProviders = authProviders.filter(
    (p) => !providerStatuses[p.id]?.enabled && !providerStatuses[p.id]?.configured
  );
  const firstAvailableProvider = availableProviders?.length ? availableProviders[0] : null;

  return (
    <Page navId="authentication">
      <Page.Contents isLoading={isLoading}>
        <h3 className={styles.sectionHeader}>Configured authentication</h3>
        {!!enabledProviders?.length && (
          <div className={styles.cardsContainer}>
            {enabledProviders.map((provider) => (
              <ProviderCard
                key={provider.id}
                providerId={provider.id}
                displayName={provider.displayName}
                authType={provider.type}
                enabled={providerStatuses[provider.id]?.enabled}
                configPath={provider.configPath}
              />
            ))}
          </div>
        )}
        {!enabledProviders?.length && firstAvailableProvider && !isEmpty(providerStatuses) && (
          <ConfigureAuthCTA
            title={`You have no ${firstAvailableProvider.type} configuration created at the moment`}
            buttonIcon="plus-circle"
            buttonLink={getProviderUrl(firstAvailableProvider)}
            buttonTitle={`Configure ${firstAvailableProvider.type}`}
            description={`Important: if you have ${firstAvailableProvider.type} configuration enabled via the .ini file Grafana is using it.
              Configuring ${firstAvailableProvider.type} via UI will take precedence over any configuration in the .ini file.
              No changes will be written into .ini file.`}
          />
        )}
        {!!configuresProviders?.length && (
          <div className={styles.cardsContainer}>
            {configuresProviders.map((provider) => (
              <ProviderCard
                key={provider.id}
                providerId={provider.id}
                displayName={provider.displayName}
                authType={provider.protocol}
                enabled={providerStatuses[provider.id]?.enabled}
                configPath={provider.configPath}
              />
            ))}
          </div>
        )}
      </Page.Contents>
    </Page>
  );
};

const getStyles = (theme: GrafanaTheme2) => {
  return {
    cardsContainer: css`
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(288px, 1fr));
      gap: ${theme.spacing(3)};
      margin-bottom: ${theme.spacing(3)};
      margin-top: ${theme.spacing(2)};
    `,
    sectionHeader: css`
      margin-bottom: ${theme.spacing(3)};
    `,
    settingsSection: css`
      margin-top: ${theme.spacing(4)};
    `,
    settingName: css`
      padding-left: 25px;
    `,
    settingValue: css`
      white-space: break-spaces;
    `,
  };
};

export default connector(AuthConfigPageUnconnected);
