import { Card, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { ApplicationDeployment } from '@/redux/types/applications';
import React from 'react';

function DeploymentsList({ deployments }: { deployments?: ApplicationDeployment[] }) {
  return (
    <div className=''>
      {
        deployments?.map((deployment) => {
          return (
            <Card key={deployment.id}>
              <CardHeader>
                <CardTitle>#{deployment.id.slice(0, 6)}</CardTitle>
                <CardDescription></CardDescription>
              </CardHeader>
            </Card>
          )
        })
      }
    </div>
  )
}

export default DeploymentsList;
