'use client';

import React from 'react';
import DashboardPageHeader from '@/components/layout/dashboard-page-header';
import useGeneralSettings from '../hooks/use-general-settings';
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import AvatarSection from './components/AvatarSection';
import AccountSection from './components/AccountSection';
import SecuritySection from './components/SecuritySection';
import { useTranslation } from '@/hooks/use-translation';
import FeatureFlagsSettings from './components/FeatureFlagsSettings';

function Page() {
  const { t, isLoading } = useTranslation();
  const {
    onImageChange,
    user,
    username,
    email,
    usernameError,
    usernameSuccess,
    emailSent,
    isLoading: settingsLoading,
    handleUsernameChange,
    handlePasswordResetRequest,
    setUsername,
    setUsernameError,
    userSettings,
    isGettingUserSettings,
    isUpdatingFont,
    isUpdatingTheme,
    isUpdatingLanguage,
    isUpdatingAutoUpdate,
    handleThemeChange,
    handleLanguageChange,
    handleAutoUpdateChange,
    handleFontUpdate
  } = useGeneralSettings();

  return (
    <div className="container mx-auto py-6 space-y-8 max-w-4xl">
      <DashboardPageHeader label={t('settings.title')} description={t('settings.description')} />
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <AvatarSection onImageChange={onImageChange} user={user} />
        <div className="col-span-1 lg:col-span-2">
          <Tabs defaultValue="account" className="w-full">
            <TabsList className="grid w-full grid-cols-3">
              <TabsTrigger value="account">{t('settings.tabs.account')}</TabsTrigger>
              <TabsTrigger value="security">{t('settings.tabs.security')}</TabsTrigger>
              <TabsTrigger value="feature-flags">{t('settings.tabs.featureFlags')}</TabsTrigger>
            </TabsList>
            <AccountSection
              username={username}
              setUsername={setUsername}
              usernameError={usernameError}
              usernameSuccess={usernameSuccess}
              setUsernameError={setUsernameError}
              email={email}
              isLoading={settingsLoading}
              handleUsernameChange={handleUsernameChange}
              user={user}
              userSettings={
                userSettings || {
                  id: '0',
                  user_id: '0',
                  font_family: 'outfit',
                  font_size: 16,
                  language: 'en',
                  theme: 'light',
                  auto_update: true,
                  created_at: new Date().toISOString(),
                  updated_at: new Date().toISOString()
                }
              }
              isGettingUserSettings={isGettingUserSettings}
              isUpdatingFont={isUpdatingFont}
              isUpdatingTheme={isUpdatingTheme}
              isUpdatingLanguage={isUpdatingLanguage}
              isUpdatingAutoUpdate={isUpdatingAutoUpdate}
              handleThemeChange={handleThemeChange}
              handleLanguageChange={handleLanguageChange}
              handleAutoUpdateChange={handleAutoUpdateChange}
              handleFontUpdate={handleFontUpdate}
            />
            <SecuritySection
              emailSent={emailSent}
              isLoading={settingsLoading}
              handlePasswordResetRequest={handlePasswordResetRequest}
            />
            <FeatureFlagsSettings />
          </Tabs>
        </div>
      </div>
    </div>
  );
}

export default Page;
