'use client';
import React from 'react';
import DashboardPageHeader from '@/components/layout/dashboard-page-header';
import NotificationPreferencesTab from './components/preferenceTab';
import NotificationChannelsTab from './components/channelTab';
import useNotificationSettings from '../hooks/use-notification-settings';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useAppSelector } from '@/redux/hooks';
import { hasPermission } from '@/lib/permission';
import { toast } from 'sonner';

export type NotificationChannelConfig = {
  [key: string]: string;
};

const Page: React.FC = () => {
  const user = useAppSelector((state) => state.auth.user);
  const activeOrg = useAppSelector((state) => state.user.activeOrganization);
  const { smtpConfigs, isLoading, handleOnSave, preferences, handleUpdatePreference } =
    useNotificationSettings();

  const canRead = hasPermission(user, 'notification', 'read', activeOrg?.id);
  const canUpdate = hasPermission(user, 'notification', 'update', activeOrg?.id);
  const canCreate = hasPermission(user, 'notification', 'create', activeOrg?.id);

  if (!canRead) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold">Access Denied</h2>
          <p className="text-muted-foreground">
            You don't have permission to view notification settings
          </p>
        </div>
      </div>
    );
  }

  const handleSave = (data: Record<string, string>) => {
    if (smtpConfigs) {
      if (canUpdate) {
        handleOnSave(data);
      } else {
        toast.error('You do not have permission to update notification settings');
      }
    } else {
      if (canCreate) {
        handleOnSave(data);
      } else {
        toast.error('You do not have permission to create notification settings');
      }
    }
  };

  const handleUpdate = (id: string, enabled: boolean) => {
    handleUpdatePreference(id, enabled);
  };

  const showChannelsTab = canCreate || (smtpConfigs && canUpdate);

  if (!showChannelsTab) {
    return (
      <div className="container mx-auto py-6 space-y-8 max-w-4xl">
        <DashboardPageHeader
          label="Notifications"
          description="Manage your notification preferences and channels"
        />
        <NotificationPreferencesTab
          activityPreferences={preferences?.activity}
          securityPreferences={preferences?.security}
          updatePreferences={preferences?.update}
          onUpdatePreference={handleUpdate}
        />
      </div>
    );
  }

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
            handleOnSave={handleSave}
          />
        </TabsContent>
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
  );
};

export default Page;
