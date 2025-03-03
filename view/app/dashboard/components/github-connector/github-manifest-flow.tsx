import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Github, RefreshCw, Loader2 } from 'lucide-react';
import { GitHubAppCredentials, GitHubAppManifest, GitHubAppProps, GitHubAppStatus } from '@/redux/types/github';

const adjectives = [
    'cosmic',
    'quantum',
    'stellar',
    'neural',
    'cyber',
    'atomic',
    'digital',
    'nebula',
    'phoenix',
    'shadow',
] as const;

const nouns = [
    'nexus',
    'pulse',
    'matrix',
    'cipher',
    'orbit',
    'nova',
    'core',
    'forge',
    'prism',
    'vector',
] as const;

const generateRandomName = (): string => {
    const adjective = adjectives[Math.floor(Math.random() * adjectives.length)];
    const noun = nouns[Math.floor(Math.random() * nouns.length)];
    return `${adjective}-${noun}`;
};

const GitHubAppManifestComponent: React.FC<GitHubAppProps> = ({
    organization,
    webhookUrl = process.env.NEXT_PUBLIC_WEBHOOK_URL,
    appUrl = process.env.NEXT_PUBLIC_APP_URL,
    redirectUrl = process.env.NEXT_PUBLIC_REDIRECT_URL,
    onSuccess,
    onError,
}) => {
    const [appName, setAppName] = useState<string>(generateRandomName());
    const [status, setStatus] = useState<GitHubAppStatus>('initial');
    const [error, setError] = useState<string | null>(null);
    // const [registerGithubApp, { isLoading, error: registerGithubAppError }] =
    //     useRegisterGithubAppMutation();

    useEffect(() => {
        const params = new URLSearchParams(window.location.search);
        const code = params.get('code');
        const stateParam = params.get('state');
        if (code) {
            handleGitHubCallback(code, stateParam);
        }
    }, []);

    const generateState = (): string => {
        return crypto
            .getRandomValues(new Uint8Array(16))
            .reduce((acc, val) => acc + val.toString(16).padStart(2, '0'), '');
    };

    const createManifestForm = (): void => {
        const state = generateState();
        const manifest: GitHubAppManifest = {
            name: appName,
            url: appUrl || window.location.origin,
            hook_attributes: {
                url: webhookUrl || `${window.location.origin}/github/webhook`,
                active: true,
            },
            redirect_url: redirectUrl || `${window.location.origin}/dashboard`,
            callback_urls: [redirectUrl || `${window.location.origin}/dashboard`],
            public: true,
            default_permissions: {
                contents: 'read',
                issues: 'write',
                metadata: 'read',
                pull_requests: 'write',
            },
            default_events: ['issues', 'issue_comment', 'pull_request', 'push'],
            setup_url: `${window.location.origin}/dashboard`,
            setup_on_update: true,
        };

        const form = document.createElement('form');
        form.method = 'post';
        form.action = organization
            ? `https://github.com/organizations/${organization}/settings/apps/new?state=${state}`
            : `https://github.com/settings/apps/new?state=${state}`;

        const input = document.createElement('input');
        input.type = 'hidden';
        input.name = 'manifest';
        input.value = JSON.stringify(manifest);
        form.appendChild(input);

        document.body.appendChild(form);
        form.submit();
        document.body.removeChild(form);

        setStatus('redirecting');
    };

    const handleGitHubCallback = async (code: string, stateParam: string | null): Promise<void> => {
        setStatus('converting');
        try {
            const response = await fetch(
                `https://api.github.com/app-manifests/${code}/conversions`,
                {
                    method: 'POST',
                    headers: {
                        Accept: 'application/vnd.github.v3+json',
                    },
                },
            );

            if (!response.ok) {
                console.log('response', response);
                throw new Error('Failed to convert manifest');
            }

            const credentials: GitHubAppCredentials = await response.json();

            // registerGithubApp({
            //     appId: credentials.id.toString(),
            //     appSlug: credentials.slug,
            //     clientId: credentials.client_id,
            //     clientSecret: credentials.client_secret,
            //     privateKey: credentials.pem,
            //     webhookSecret: credentials.webhook_secret,
            // });

            setStatus('success');
            onSuccess?.(credentials);
        } catch (err) {
            const errorMessage = err instanceof Error ? err.message : 'An unknown error occurred';
            setError(errorMessage);
            setStatus('error');
            onError?.(err instanceof Error ? err : new Error(errorMessage));
        }
    };

    return (
        <Card className="w-[400px]">
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    <Github size={24} />
                    Create GitHub App
                </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
                {status === 'initial' && (
                    <>
                        <div className="space-y-2">
                            <label className="text-sm font-medium">App Name</label>
                            <div className="flex gap-2">
                                <Input
                                    value={appName}
                                    onChange={(e) => setAppName(e.target.value)}
                                    placeholder="Enter app name"
                                    className="flex-1"
                                />
                                <Button
                                    variant="outline"
                                    size="icon"
                                    onClick={() => setAppName(generateRandomName())}
                                    title="Generate new name"
                                >
                                    <RefreshCw size={16} />
                                </Button>
                            </div>
                        </div>
                        {error && (
                            <Alert variant="destructive">
                                <AlertDescription>{error}</AlertDescription>
                            </Alert>
                        )}

                        <Button className="w-full" onClick={createManifestForm}>
                            Create GitHub App
                        </Button>
                    </>
                )}

                {(status === 'redirecting' || status === 'converting') && (
                    <div className="flex flex-col items-center gap-4 py-8">
                        <Loader2 className="h-8 w-8 animate-spin" />
                        <p>
                            {status === 'redirecting'
                                ? 'Redirecting to GitHub...'
                                : 'Setting up your GitHub App...'}
                        </p>
                    </div>
                )}

                {status === 'success' && (
                    <Alert>
                        <AlertDescription className="text-green-600">
                            GitHub App created successfully!
                        </AlertDescription>
                    </Alert>
                )}
            </CardContent>
        </Card>
    );
};

export default GitHubAppManifestComponent;
