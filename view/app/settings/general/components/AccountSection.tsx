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
      setVerificationError('Failed to send verification email. Please try again later.');
    }
  };

  return (
    <TabsContent value="account" className="space-y-4 mt-4">
      <Card>
        <CardHeader>
          <CardTitle>Account Information</CardTitle>
          <CardDescription>Update your account details</CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="space-y-2">
            <Label htmlFor="username" className="flex items-center gap-2">
              <User size={16} />
              Username
            </Label>
            <div className="flex gap-2">
              <Input
                id="username"
                value={username}
                onChange={(e) => {
                  setUsername(e.target.value);
                  setUsernameError('');
                }}
                placeholder="Enter your username"
              />
              <Button
                onClick={handleUsernameChange}
                disabled={isLoading || username === user.username}
              >
                Update
              </Button>
            </div>

            {usernameError && <p className="text-sm text-red-500">{usernameError}</p>}

            {usernameSuccess && (
              <Alert variant="default">
                <CheckCircle className="h-4 w-4" />
                <AlertTitle>Success</AlertTitle>
                <AlertDescription>Your username has been updated successfully!</AlertDescription>
              </Alert>
            )}
          </div>

          <Separator />

          <div className="space-y-2">
            <Label htmlFor="email" className="flex items-center gap-2">
              <Mail size={16} />
              Email Address
            </Label>
            <div className="flex flex-col gap-2">
              <Input id="email" value={email} readOnly disabled className="bg-muted/50" />
              {!user.is_verified && (
                <div className="space-y-2">
                  <Alert variant="destructive">
                    <AlertCircle className="h-4 w-4" />
                    <AlertTitle>Email Not Verified</AlertTitle>
                    <AlertDescription>
                      Your email address is not verified. Please verify your email to access all
                      features.
                    </AlertDescription>
                  </Alert>
                  <Button
                    onClick={handleSendVerification}
                    disabled={isSendingVerification || verificationSent}
                    variant="outline"
                    className="w-full"
                  >
                    {isSendingVerification
                      ? 'Sending...'
                      : verificationSent
                        ? 'Verification Email Sent'
                        : 'Send Verification Email'}
                  </Button>
                  {verificationError && <p className="text-sm text-red-500">{verificationError}</p>}
                  {verificationSent && (
                    <Alert variant="default">
                      <CheckCircle className="h-4 w-4" />
                      <AlertTitle>Verification Email Sent</AlertTitle>
                      <AlertDescription>
                        Please check your email and click the verification link to verify your email
                        address.
                      </AlertDescription>
                    </Alert>
                  )}
                </div>
              )}
              <p className="text-sm text-muted-foreground">
                Contact support to change your email address
              </p>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Account Preferences</CardTitle>
          <CardDescription>Manage your notification and display settings</CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground text-sm flex justify-between items-center">
            <span>Appearance</span> <ModeToggle />
          </p>
        </CardContent>
      </Card>
    </TabsContent>
  );
}

export default AccountSection;
