'use client';
import React from 'react';
import { Button } from '@/components/ui/button';
import { Form } from '@/components/ui/form';
import FormInputField from '@/components/ui/form-input-field';
import { FormSelectTagInputField } from '@/components/ui/form-select-tag-field';
import { BuildPack, Environment } from '@/redux/types/deploy-form';
import useUpdateDeployment from '../../hooks/use_update_deployment';
import { parsePort } from '../../utils/parsePort';
import { useAppSelector } from '@/redux/hooks';
import { hasPermission } from '@/lib/permission';

interface DeployConfigureProps {
  application_name?: string;
  environment?: Environment;
  branch?: string;
  port?: string;
  domain?: string;
  repository?: string;
  build_pack?: BuildPack;
  env_variables?: Record<string, string>;
  build_variables?: Record<string, string>;
  pre_run_commands?: string;
  post_run_commands?: string;
  application_id?: string;
  dockerFilePath?: string;
  base_path?: string;
}

export const DeployConfigureForm = ({
  application_name = '',
  environment = Environment.Production,
  branch = '',
  port = '3000',
  domain = '',
  repository = '',
  build_pack = BuildPack.Dockerfile,
  env_variables = {},
  build_variables = {},
  pre_run_commands = '',
  post_run_commands = '',
  application_id = '',
  dockerFilePath = '/Dockerfile',
  base_path = '/'
}: DeployConfigureProps) => {
  const user = useAppSelector((state) => state.auth.user);
  const activeOrg = useAppSelector((state) => state.user.activeOrganization);
  const canUpdate = hasPermission(user, 'deploy', 'update', activeOrg?.id);

  const { validateEnvVar, form, onSubmit, isLoading, domains } = useUpdateDeployment({
    name: application_name,
    pre_run_command: pre_run_commands,
    post_run_command: post_run_commands,
    build_variables,
    environment_variables: env_variables,
    port: parsePort(port) || 3000,
    force: true,
    id: application_id,
    DockerfilePath: dockerFilePath,
    base_path
  });

  const renderReadOnlyField = (label: string, value: string | undefined, description: string) => (
    <div className="space-y-2">
      <label className="text-sm font-medium">{label}</label>
      <div className="px-3 py-2 border rounded-md bg-muted text-muted-foreground">
        {value || '-'}
      </div>
      <p className="text-sm text-muted-foreground">{description}</p>
    </div>
  );

  if (!canUpdate) {
    return (
      <div className="space-y-8">
        <div className="grid sm:grid-cols-2 gap-4">
          {renderReadOnlyField('Application Name', application_name, 'Application name')}
          {renderReadOnlyField('Port', port, 'Port on which your application will be available')}
          {renderReadOnlyField(
            'Base Path',
            base_path,
            'The build context path for your application'
          )}
          {renderReadOnlyField(
            'Dockerfile Path',
            dockerFilePath,
            'Path of the dockerfile relative to the base path'
          )}
        </div>

        <div className="grid sm:grid-cols-2 gap-4">
          {renderReadOnlyField(
            'Environment Variables',
            Object.entries(env_variables)
              .map(([key, value]) => `${key}=${value}`)
              .join('\n'),
            'Environment variables for the application'
          )}
          {renderReadOnlyField(
            'Build Variables',
            Object.entries(build_variables)
              .map(([key, value]) => `${key}=${value}`)
              .join('\n'),
            'Build variables for the application'
          )}
        </div>

        <div className="grid sm:grid-cols-2 gap-4">
          {renderReadOnlyField(
            'Pre Run Commands',
            pre_run_commands,
            'Commands to run before deployment'
          )}
          {renderReadOnlyField(
            'Post Run Commands',
            post_run_commands,
            'Commands to run after deployment'
          )}
        </div>

        <div className="grid sm:grid-cols-2 gap-4">
          {renderReadOnlyField('Environment', environment, 'Environment of the deployment')}
          {renderReadOnlyField('Branch', branch, 'Branch from where you want to deploy')}
        </div>

        <div className="grid sm:grid-cols-2 gap-4">
          {renderReadOnlyField('Domain', domain, 'Domain on which your application is available')}
          {renderReadOnlyField('Build Pack', build_pack, 'Build pack used for the application')}
        </div>
      </div>
    );
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        <div className="grid sm:grid-cols-2 gap-4">
          <FormInputField
            form={form}
            label="Application Name"
            name="name"
            description="Application name"
            placeholder="Application name"
          />
          <FormInputField
            form={form}
            label="Port"
            name="port"
            description="Port on which your application will be available"
            placeholder="3000"
            validator={(value) => parsePort(value) !== null}
          />
          <FormInputField
            form={form}
            label="Base Path"
            name="base_path"
            description="The build context path for your application"
            placeholder="/"
            required={false}
          />
          <FormInputField
            form={form}
            label="Dockerfile Path"
            name="DockerfilePath"
            description="Path of the dockerfile relative to the base path"
            placeholder="Dockerfile"
            required={false}
          />
        </div>

        <div className="grid sm:grid-cols-2 gap-4">
          <FormSelectTagInputField
            form={form}
            label="Environment Variables"
            name="environment_variables"
            description="Type KEY=VALUE and press Enter to add"
            placeholder="NODE_ENV=production"
            required={false}
            validator={validateEnvVar}
            defaultValues={env_variables}
          />
          <FormSelectTagInputField
            form={form}
            label="Build Variables"
            name="build_variables"
            description="Type KEY=VALUE and press Enter to add"
            placeholder="DEBUG=false"
            required={false}
            validator={validateEnvVar}
            defaultValues={build_variables}
          />
        </div>

        <div className="grid sm:grid-cols-2 gap-4">
          <FormInputField
            form={form}
            label="Pre Run Commands"
            name="pre_run_command"
            description="Commands to run before deployment"
            placeholder="npm install"
            required={false}
          />
          <FormInputField
            form={form}
            label="Post Run Commands"
            name="post_run_command"
            description="Commands to run after deployment"
            placeholder="npm run test"
            required={false}
          />
        </div>

        <div className="grid sm:grid-cols-2 gap-4">
          {renderReadOnlyField('Environment', environment, 'Environment of the deployment')}
          {renderReadOnlyField('Branch', branch, 'Branch from where you want to deploy')}
        </div>

        <div className="grid sm:grid-cols-2 gap-4">
          {renderReadOnlyField('Domain', domain, 'Domain on which your application is available')}
          {renderReadOnlyField('Build Pack', build_pack, 'Build pack used for the application')}
        </div>

        <Button type="submit" className="w-full cursor-pointer" disabled={isLoading}>
          {isLoading ? 'Updating...' : 'Update Deployment'}
        </Button>
      </form>
    </Form>
  );
};
