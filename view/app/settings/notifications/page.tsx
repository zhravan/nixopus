'use client';
import React from 'react';
import DashboardPageHeader from '@/components/layout/dashboard-page-header';
import NotificationPreferencesTab from './components/preferenceTab';
import NotificationChannelsTab from './components/channelTab';
import useNotificationSettings from '../hooks/use-notification-settings';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';

export type NotificationChannelConfig = {
  [key: string]: string;
};

const Page: React.FC = () => {
  const { smtpConfigs, isLoading, handleOnSave, preferences, handleUpdatePreference } =
    useNotificationSettings();

  return (
    <div className="container mx-auto py-6 space-y-8 max-w-4xl">
      <DashboardPageHeader
        label="Notifications"
        description="Manage your notification preferences and channels"
      />
      <Tabs defaultValue="channels" className="w-full">
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="channels">Channels</TabsTrigger>
          <TabsTrigger value="preferences">Preferences</TabsTrigger>
        </TabsList>
        <TabsContent value="channels">
          <NotificationChannelsTab
            smtpConfigs={smtpConfigs}
            isLoading={isLoading}
            handleOnSave={handleOnSave}
          />
        </TabsContent>
        <TabsContent value="preferences">
          <NotificationPreferencesTab
            activityPreferences={preferences?.activity}
            securityPreferences={preferences?.security}
            updatePreferences={preferences?.update}
            onUpdatePreference={handleUpdatePreference}
          />
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default Page;
