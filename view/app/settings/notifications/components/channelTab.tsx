'use client';
import React from 'react';
import { Mail } from 'lucide-react';
import NotificationChannelCard from './channel';
import { SMTPConfig } from '@/redux/types/notification';
import { useTranslation } from '@/hooks/use-translation';

interface NotificationChannelsProps {
  smtpConfigs: SMTPConfig | undefined;
  isLoading: boolean;
  handleOnSave: (data: Record<string, string>) => void;
}

export const NotificationChannels: React.FC<NotificationChannelsProps> = ({
  smtpConfigs,
  isLoading,
  handleOnSave
}) => {
  const { t } = useTranslation();

  return (
    <div className="grid gap-6 md:grid-cols-1">
      <NotificationChannelCard
        title={t('settings.notifications.channels.email.title')}
        description={t('settings.notifications.channels.email.description')}
        icon={<Mail className="h-5 w-5 text-primary" />}
        configData={{
          smtpServer: smtpConfigs?.host || '',
          port: smtpConfigs?.port.toString() || '587',
          username: smtpConfigs?.username || '',
          password: smtpConfigs?.password || '',
          fromEmail: smtpConfigs?.from_email || '',
          fromName: smtpConfigs?.from_name || ''
        }}
        onSave={(data) => handleOnSave(data)}
        isLoading={isLoading}
        channelType="email"
      />
    </div>
  );
};

export default NotificationChannels;
