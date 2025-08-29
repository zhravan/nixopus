'use client';
import React from 'react';
import DashboardPageHeader from '@/components/layout/dashboard-page-header';
import NotificationPreferencesTab from './components/preferenceTab';
import NotificationChannelsTab from './components/channelTab';
import useNotificationSettings from '../hooks/use-notification-settings';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useTranslation } from '@/hooks/use-translation';
import { SMTPFormData } from '@/redux/types/notification';
import { useFeatureFlags } from '@/hooks/features_provider';
import Skeleton from '@/app/file-manager/components/skeleton/Skeleton';
import DisabledFeature from '@/components/features/disabled-feature';
import { FeatureNames } from '@/types/feature-flags';
import { ResourceGuard } from '@/components/rbac/PermissionGuard';
import { useRBAC } from '@/lib/rbac';

export type NotificationChannelConfig = {
  [key: string]: string;
};

const Page: React.FC = () => {
  const { t } = useTranslation();
  const {
    smtpConfigs,
    slackConfig,
    discordConfig,
    isLoading,
    handleOnSave,
    handleCreateWebhookConfig,
    handleUpdateWebhookConfig,
    handleDeleteWebhookConfig,
    preferences,
    handleUpdatePreference
  } = useNotificationSettings();
  const { isFeatureEnabled, isLoading: isFeatureFlagsLoading } = useFeatureFlags();
  const { canAccessResource } = useRBAC();
  const hasFeatureFlagsReadPermission = canAccessResource('feature-flags', 'read');

  if (isFeatureFlagsLoading) {
    return <Skeleton />;
  }

  if (!isFeatureEnabled(FeatureNames.FeatureNotifications)) {
    return <DisabledFeature />;
  }

  const handleSave = (data: SMTPFormData) => {
    if (smtpConfigs) {
      handleOnSave(data);
    } else {
      handleOnSave(data);
    }
  };

  const handleSaveSlack = (data: Record<string, string>) => {
    if (slackConfig) {
      handleUpdateWebhookConfig({
        type: 'slack',
        webhook_url: data.webhook_url,
        is_active: data.is_active === 'true'
      });
    } else {
      handleCreateWebhookConfig({
        type: 'slack',
        webhook_url: data.webhook_url
      });
    }
  };

  const handleSaveDiscord = (data: Record<string, string>) => {
    if (discordConfig) {
      handleUpdateWebhookConfig({
        type: 'discord',
        webhook_url: data.webhook_url,
        is_active: data.is_active === 'true'
      });
    } else {
      handleCreateWebhookConfig({
        type: 'discord',
        webhook_url: data.webhook_url
      });
    }
  };

  const handleUpdate = (id: string, enabled: boolean) => {
    handleUpdatePreference(id, enabled);
  };

  return (
    <ResourceGuard resource="notification" action="read">
      <div className="container mx-auto py-6 space-y-8 max-w-4xl">
        <DashboardPageHeader
          label={t('settings.notifications.page.title')}
          description={t('settings.notifications.page.description')}
        />
        <Tabs defaultValue={hasFeatureFlagsReadPermission ? "channels" : "preferences"} className="w-full">
          <TabsList className={`grid w-full grid-cols-2`}>
            <ResourceGuard resource="notification" action="create">
              <TabsTrigger value="channels">
                {t('settings.notifications.page.tabs.channels')}
              </TabsTrigger>
            </ResourceGuard>
            <TabsTrigger value="preferences">
              {t('settings.notifications.page.tabs.preferences')}
            </TabsTrigger>
          </TabsList>
          <ResourceGuard resource="notification" action="create">
            <TabsContent value="channels">
              <NotificationChannelsTab
                smtpConfigs={smtpConfigs}
                slackConfig={slackConfig}
                discordConfig={discordConfig}
                isLoading={isLoading}
                handleOnSave={handleSave}
                handleOnSaveSlack={handleSaveSlack}
                handleOnSaveDiscord={handleSaveDiscord}
              />
            </TabsContent>
          </ResourceGuard>
          <TabsContent value="preferences">
            <NotificationPreferencesTab
              activityPreferences={preferences?.activity}
              securityPreferences={preferences?.security}
              updatePreferences={preferences?.update}
              onUpdatePreference={handleUpdate}
            />
          </TabsContent>
        </Tabs>
      </div>
    </ResourceGuard>
  );
};

export default Page;
