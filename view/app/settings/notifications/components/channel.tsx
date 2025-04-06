'use client';
import React, { useState, useEffect, ReactNode } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { useTranslation } from '@/hooks/use-translation';

export interface NotificationChannelField {
  id: string;
  label: string;
  placeholder: string;
  type?: string;
  required: boolean;
}

export interface NotificationChannelProps {
  title: string;
  description: string;
  icon: ReactNode;
  configData?: Record<string, string>;
  onSave?: (data: Record<string, string>) => void;
  isLoading?: boolean;
  channelType: 'email' | 'slack' | 'webhook';
}

const NotificationChannelCard: React.FC<NotificationChannelProps> = ({
  title,
  description,
  icon,
  configData = {},
  onSave,
  isLoading,
  channelType
}) => {
  const { t } = useTranslation();
  const [formData, setFormData] = useState<Record<string, string>>(configData);
  const [isFormValid, setIsFormValid] = useState<boolean>(false);

  useEffect(() => {
    setFormData(configData);
  }, [configData]);

  const getChannelFields = (): NotificationChannelField[] => {
    switch (channelType) {
      case 'email':
        return [
          {
            id: 'smtpServer',
            label: t('settings.notifications.channels.email.fields.smtpServer.label'),
            placeholder: t('settings.notifications.channels.email.fields.smtpServer.placeholder'),
            required: true
          },
          {
            id: 'port',
            label: t('settings.notifications.channels.email.fields.port.label'),
            placeholder: t('settings.notifications.channels.email.fields.port.placeholder'),
            required: true
          },
          {
            id: 'username',
            label: t('settings.notifications.channels.email.fields.username.label'),
            placeholder: t('settings.notifications.channels.email.fields.username.placeholder'),
            required: true
          },
          {
            id: 'password',
            label: t('settings.notifications.channels.email.fields.password.label'),
            placeholder: t('settings.notifications.channels.email.fields.password.placeholder'),
            type: 'password',
            required: true
          },
          {
            id: 'fromEmail',
            label: t('settings.notifications.channels.email.fields.fromEmail.label'),
            placeholder: t('settings.notifications.channels.email.fields.fromEmail.placeholder'),
            required: true
          },
          {
            id: 'fromName',
            label: t('settings.notifications.channels.email.fields.fromName.label'),
            placeholder: t('settings.notifications.channels.email.fields.fromName.placeholder'),
            required: true
          }
        ];
      default:
        return [];
    }
  };

  const fields = getChannelFields();

  useEffect(() => {
    const valid = fields
      .filter((field) => field.required)
      .every((field) => formData[field.id] && formData[field.id].trim() !== '');
    setIsFormValid(valid);
  }, [formData, fields]);

  const handleInputChange = (id: string, value: string) => {
    setFormData((prev) => ({
      ...prev,
      [id]: value
    }));
  };

  const handleSave = () => {
    if (isFormValid) {
      onSave && onSave(formData);
    }
  };

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <div className="flex items-center space-x-3">
          {icon}
          <div>
            <CardTitle className="text-lg">{title}</CardTitle>
            <CardDescription>{description}</CardDescription>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-4 pt-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {fields.map((field) => (
              <div className="space-y-2" key={field.id}>
                <Label htmlFor={`${channelType}-${field.id}`}>
                  {field.label} {field.required && <span className="text-red-500">*</span>}
                </Label>
                <Input
                  id={`${channelType}-${field.id}`}
                  type={field.type || 'text'}
                  value={formData[field.id] || ''}
                  onChange={(e) => handleInputChange(field.id, e.target.value)}
                  placeholder={field.placeholder}
                />
              </div>
            ))}
          </div>

          <div className="pt-2 flex space-x-2 justify-end">
            <Button onClick={handleSave} disabled={!isFormValid || isLoading}>
              {t('settings.notifications.channels.email.buttons.save')}
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};

export default NotificationChannelCard;
