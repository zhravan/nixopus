'use client';
import React from 'react';
import { Button } from '@/components/ui/button';
import { ExternalLink, MoreVertical, RotateCcw, TrashIcon } from 'lucide-react';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu';
import { Tooltip, TooltipProvider, TooltipTrigger, TooltipContent } from '@/components/ui/tooltip';
import { Application } from '@/redux/types/applications';
import { DeleteDialog } from '@/components/ui/delete-dialog';
import {
  useRedeployApplicationMutation,
  useRestartApplicationMutation
} from '@/redux/services/deploy/applicationsApi';
import { useDeleteApplicationMutation } from '@/redux/services/deploy/applicationsApi';
import { useRouter } from 'next/navigation';
import { toast } from 'sonner';

const ApplicationDetailsHeader = ({ application }: { application?: Application }) => {
  const [redeployApplication, { isLoading: isRedeploying }] = useRedeployApplicationMutation();
  const [deleteApplication, { isLoading: isDeleting }] = useDeleteApplicationMutation();
  const router = useRouter();
  const [restartApplication, { isLoading: isRestarting }] = useRestartApplicationMutation();

  const handleDelete = async () => {
    try {
      await deleteApplication({
        id: application?.id || ''
      }).unwrap();
      toast.success('Application deleted successfully');
      router.push('/self-host');
    } catch (error) {
      toast.error('Failed to delete application');
    }
  };

  const handleRestart = async () => {
    try {
      await restartApplication({ id: application?.deployments?.[0]?.id || '' }).unwrap();
      toast.success('Application restart started');
    } catch (error) {
      toast.error('Failed to restart application');
    }
  };

  const handleRedeploy = async (forceWithoutCache: boolean) => {
    try {
      await redeployApplication({
        id: application?.id || '',
        force: true,
        force_without_cache: forceWithoutCache
      }).unwrap();
      router.push('/self-host/application/' + application?.id + '?logs=true');
      toast.success('Application redeployment started');
    } catch (error) {
      toast.error('Failed to redeploy application');
    }
  };

  return (
    <div className="flex flex-col items-start justify-between gap-4 sm:flex-row sm:items-center">
      <div className="flex items-start">
        <div className="mr-2">
          <div className="flex items-center gap-2">
            <h1 className="text-2xl font-bold capitalize sm:text-3xl">{application?.name}</h1>
            <Button
              variant="ghost"
              size="icon"
              onClick={() => window.open('https://' + application?.domain, '_blank')}
              aria-label="Open application in new tab"
            >
              <ExternalLink className="h-5 w-5" />
            </Button>
          </div>
        </div>
      </div>
      <div className="flex items-center gap-2">
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant="secondary"
                size="icon"
                disabled={isRestarting}
                onClick={handleRestart}
              >
                <RotateCcw className="h-4 w-4" />
              </Button>
            </TooltipTrigger>
            <TooltipContent>
              <p>Restart Application</p>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
        <DeleteDialog
          title={`Delete ${application?.name}`}
          description={`Are you sure you want to delete ${application?.name}? This action cannot be undone.`}
          onConfirm={handleDelete}
          trigger={
            <Button variant="outline" size="icon">
              <TrashIcon className="h-4 w-4" />
            </Button>
          }
          isDeleting={isDeleting}
          variant="destructive"
          icon={TrashIcon}
        />
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline" size="icon">
              <MoreVertical className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={() => handleRedeploy(true)} disabled={isRedeploying}>
              Re-Deploy
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => handleRedeploy(false)} disabled={isRedeploying}>
              Force Deploy Without Cache
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  );
};

export default ApplicationDetailsHeader;
