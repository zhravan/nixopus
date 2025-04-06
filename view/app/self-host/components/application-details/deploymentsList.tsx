import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useRollbackApplicationMutation } from '@/redux/services/deploy/applicationsApi';
import { ApplicationDeployment } from '@/redux/types/applications';
import { Undo, Eye, Clock } from 'lucide-react';
import { useRouter } from 'next/navigation';
import React from 'react';
import { useTranslation } from '@/hooks/use-translation';

function DeploymentsList({ deployments }: { deployments?: ApplicationDeployment[] }) {
  const { t } = useTranslation();
  const router = useRouter();
  const [rollBackApplication, { isLoading }] = useRollbackApplicationMutation();

  const formatDate = (created_at: string) =>
    deployments
      ? new Date(created_at).toLocaleString('en-US', {
          day: 'numeric',
          month: 'short',
          year: 'numeric',
          hour: '2-digit',
          minute: '2-digit'
        })
      : 'N/A';

  const calculateRunTime = (updated_at: string, created_at: string) => {
    const createdTime = new Date(created_at).getTime();
    const updatedTime = new Date(updated_at).getTime();
    const duration = updatedTime - createdTime;

    const minutes = Math.floor(duration / 60000);
    const seconds = Math.floor((duration % 60000) / 1000);

    return `${minutes}m ${seconds}s`;
  };

  return (
    <div className="space-y-6">
      {deployments && deployments.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {deployments.map((deployment) => (
            <Card key={deployment.id} className="w-full hover:shadow-md transition-shadow">
              <CardHeader className="pb-2">
                <div className="flex justify-between items-center">
                  <div>
                    <CardTitle className="text-lg font-medium">
                      {deployment.container_name?.startsWith('/')
                        ? deployment.container_name.slice(1)
                        : deployment.container_name}
                    </CardTitle>
                    <CardDescription>
                      <span className="text-xs">
                        {deployment.status?.status ||
                          t('selfHost.deployments.list.card.status.completed')}
                      </span>
                    </CardDescription>
                  </div>
                  <Button
                    size="icon"
                    onClick={() => rollBackApplication({ id: deployment.id })}
                    disabled={isLoading}
                    title={t('selfHost.deployments.list.card.rollback.title')}
                  >
                    <Undo className="h-4 w-4" />
                  </Button>
                </div>
                <CardDescription>{formatDate(deployment.created_at)}</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="flex flex-col space-y-1">
                  <div className="flex justify-between text-sm">
                    <span className="flex items-center gap-2">
                      <Clock className="h-4 w-4" /> {t('selfHost.deployments.list.card.runTime')}:
                    </span>
                    <span className="font-medium">
                      {calculateRunTime(
                        deployment.status?.updated_at as string,
                        deployment.status?.created_at as string
                      )}
                    </span>
                  </div>
                </div>
                <Button
                  className="w-full mt-3 cursor-pointer"
                  variant="secondary"
                  onClick={() => {
                    router.push(
                      `/self-host/application/${deployment.application_id}/deployments/${deployment.id}`
                    );
                  }}
                >
                  <Eye className="h-4 w-4 mr-2" /> {t('selfHost.deployments.list.card.viewButton')}
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>
      ) : (
        <div className="text-center py-8 rounded-lg border border-gray-200">
          <p>{t('selfHost.deployments.list.noDeployments')}</p>
        </div>
      )}
    </div>
  );
}

export default DeploymentsList;
