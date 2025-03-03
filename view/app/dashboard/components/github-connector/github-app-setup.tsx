"use client";
import React, { useState } from 'react';
import { Card, CardContent } from '@/components/ui/card';
import { GitHubAppCredentials } from '@/redux/types/github';
import GitHubAppManifestComponent from './github-manifest-flow';
import GithubInstaller from './github-app-installer';

interface GitHubAppSetupProps {
    organization?: string;
    onSetupComplete?: (credentials: GitHubAppCredentials, installationId: string) => void;
}

const GitHubAppSetup: React.FC<GitHubAppSetupProps> = ({ organization, onSetupComplete }) => {
    const [setupStage, setSetupStage] = useState<'registration' | 'installation'>('registration');
    const [credentials, setCredentials] = useState<GitHubAppCredentials | null>(null);
    const [error, setError] = useState<string | null>(null);

    const handleRegistrationSuccess = (creds: GitHubAppCredentials) => {
        setCredentials(creds);
        setSetupStage('installation');
    };

    const handleRegistrationError = (error: Error) => {
        setError(`Registration failed: ${error.message}`);
    };

    const handleInstallationSuccess = (installationId: string) => {
        if (credentials) {
            onSetupComplete?.(credentials, installationId);
        }
    };

    const handleInstallationError = (error: Error) => {
        setError(`Installation failed: ${error.message}`);
    };

    return (
        <div className="flex flex-col items-center space-y-4 p-4">
            {setupStage === 'registration' ? (
                <GitHubAppManifestComponent
                    organization={organization}
                    onSuccess={handleRegistrationSuccess}
                    onError={handleRegistrationError}
                />
            ) : credentials ? (
                <GithubInstaller
                    appSlug={credentials.slug}
                    organization={organization}
                    callbackUrl={''}
                />
            ) : null}

            {error && (
                <Card className="w-[400px] bg-red-50">
                    <CardContent className="pt-6">
                        <p className="text-red-600">{error}</p>
                    </CardContent>
                </Card>
            )}
        </div>
    );
};

export default GitHubAppSetup;
