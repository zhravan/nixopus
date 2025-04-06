'use client';

import React, { useState } from 'react';
import { Mail, User, CheckCircle, AlertCircle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { TabsContent } from '@/components/ui/tabs';
import { Separator } from '@/components/ui/separator';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { User as UserType } from '@/redux/types/user';
import { ModeToggle } from '@/components/ui/theme-toggler';
import { useSendVerificationEmailMutation } from '@/redux/services/users/authApi';
import { useTranslation } from '@/hooks/use-translation';
import { LanguageSwitcher } from '@/components/language-switcher';

interface AccountSectionProps {
  username: string;
  setUsername: (username: string) => void;
  usernameError: string;
  usernameSuccess: boolean;
  setUsernameError: (error: string) => void;
  email: string;
  isLoading: boolean;
  handleUsernameChange: () => void;
  user: UserType;
}

function AccountSection({
  username,
  setUsername,
  usernameError,
  usernameSuccess,
  setUsernameError,
  email,
  isLoading,
  handleUsernameChange,
  user
}: AccountSectionProps) {
  const { t } = useTranslation();
  const [sendVerificationEmail, { isLoading: isSendingVerification }] =
    useSendVerificationEmailMutation();
  const [verificationSent, setVerificationSent] = useState(false);
  const [verificationError, setVerificationError] = useState('');

  const handleSendVerification = async () => {
    try {
      await sendVerificationEmail().unwrap();
      setVerificationSent(true);
      setVerificationError('');
    } catch (error) {
      setVerificationError(t('settings.account.email.notVerified.error'));
    }
  };

  return (
    <TabsContent value="account" className="space-y-4 mt-4">
      <Card>
        <CardHeader>
          <CardTitle>{t('settings.account.title')}</CardTitle>
          <CardDescription>{t('settings.account.description')}</CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="space-y-2">
            <Label htmlFor="username" className="flex items-center gap-2">
              <User size={16} />
              {t('settings.account.username.label')}
            </Label>
            <div className="flex gap-2">
              <Input
                id="username"
                value={username}
                onChange={(e) => {
                  setUsername(e.target.value);
                  setUsernameError('');
                }}
                placeholder={t('settings.account.username.placeholder')}
              />
              <Button
                onClick={handleUsernameChange}
                disabled={isLoading || username === user.username}
              >
                {t('settings.account.username.update')}
              </Button>
            </div>

            {usernameError && <p className="text-sm text-red-500">{usernameError}</p>}

            {usernameSuccess && (
              <Alert variant="default">
                <CheckCircle className="h-4 w-4" />
                <AlertTitle>Success</AlertTitle>
                <AlertDescription>{t('settings.account.username.success')}</AlertDescription>
              </Alert>
            )}
          </div>

          <Separator />

          <div className="space-y-2">
            <Label htmlFor="email" className="flex items-center gap-2">
              <Mail size={16} />
              {t('settings.account.email.label')}
            </Label>
            <div className="flex flex-col gap-2">
              <Input id="email" value={email} readOnly disabled className="bg-muted/50" />
              {!user.is_verified && (
                <div className="space-y-2">
                  <Alert variant="destructive">
                    <AlertCircle className="h-4 w-4" />
                    <AlertTitle>{t('settings.account.email.notVerified.title')}</AlertTitle>
                    <AlertDescription>
                      {t('settings.account.email.notVerified.description')}
                    </AlertDescription>
                  </Alert>
                  <Button
                    onClick={handleSendVerification}
                    disabled={isSendingVerification || verificationSent}
                    variant="outline"
                    className="w-full"
                  >
                    {isSendingVerification
                      ? t('settings.account.email.notVerified.sending')
                      : verificationSent
                        ? t('settings.account.email.notVerified.sent')
                        : t('settings.account.email.notVerified.sendButton')}
                  </Button>
                  {verificationError && <p className="text-sm text-red-500">{verificationError}</p>}
                  {verificationSent && (
                    <Alert variant="default">
                      <CheckCircle className="h-4 w-4" />
                      <AlertTitle>{t('settings.account.email.notVerified.sent')}</AlertTitle>
                      <AlertDescription>
                        {t('settings.account.email.notVerified.checkEmail')}
                      </AlertDescription>
                    </Alert>
                  )}
                </div>
              )}
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>{t('settings.account.preferences.title')}</CardTitle>
          <CardDescription>{t('settings.account.preferences.description')}</CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground text-sm flex justify-between items-center">
            <span>{t('settings.account.preferences.appearance')}</span> <ModeToggle />
          </p>
        </CardContent>
      </Card>

      <div className="mt-6">
        <Card>
          <CardHeader>
            <CardTitle>{t('settings.preferences.language.title')}</CardTitle>
            <CardDescription>{t('settings.preferences.language.description')}</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex items-center justify-between">
              <p className="text-sm text-gray-500">{t('settings.preferences.language.select')}</p>
              <LanguageSwitcher />
            </div>
          </CardContent>
        </Card>
      </div>
    </TabsContent>
  );
}

export default AccountSection;
