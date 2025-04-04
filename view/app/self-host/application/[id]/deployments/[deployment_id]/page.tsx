'use client';
import ApplicationLogs from '@/app/self-host/components/application-details/logs';
import useDeploymentDetails from '@/app/self-host/hooks/use_deployment_details';
import React from 'react';

function page() {
  const { deployment, logs } = useDeploymentDetails();

  return (
    <div className="container mx-auto py-6 space-y-8 max-w-4xl 2xl:max-w-7xl">
      <div className="mt-5 flex-col flex mb-4">
        <span className="text-2xl font-bold">Deployment Details</span>
      </div>
      <ApplicationLogs
        logs={logs}
        onRefresh={() => {}}
        currentPage={1}
        setCurrentPage={(page: number) => {}}
      />
    </div>
  );
}

export default page;
