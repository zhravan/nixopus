'use client';
import React from 'react';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import ApplicationLogs from '../../components/application-details/logs';
import Monitor from '../../components/application-details/monitor';
import Configuration from '../../components/application-details/configuration';
import DeploymentsList from '../../components/application-details/deploymentsList';
import useApplicationDetails from '../../hooks/use_application_details';

function Page() {
  const { application, currentPage, setCurrentPage } = useApplicationDetails();
  return (
    <div className="container mx-auto py-6 space-y-8 max-w-4xl mt-10">
      <Tabs defaultValue="monitoring" className="w-full">
        <TabsList>
          <TabsTrigger value="monitoring">Monitoring</TabsTrigger>
          <TabsTrigger value="configuration">Configuration</TabsTrigger>
          <TabsTrigger value="deployments">Deployments</TabsTrigger>
          <TabsTrigger value="logs">Logs</TabsTrigger>
        </TabsList>
        <TabsContent value="deployments" className="mt-6">
          <DeploymentsList deployments={application?.deployments} />
        </TabsContent>

        <TabsContent value="configuration" className="mt-6">
          <Configuration />
        </TabsContent>

        <TabsContent value="logs" className="mt-6">
          <ApplicationLogs
            logs={application?.logs}
            onRefresh={() => {}}
            currentPage={currentPage}
            setCurrentPage={setCurrentPage}
          />
        </TabsContent>

        <TabsContent value="monitoring" className="mt-6">
          <Monitor />
        </TabsContent>
      </Tabs>
    </div>
  );
}

export default Page;
