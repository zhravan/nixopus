import React from 'react';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle
} from '@/components/ui/card';
import { Button } from '@/components/ui/button';

interface TeamStatsProps {
  users: {
    id: string;
    name: string;
    role: 'Admin' | 'Member' | 'Viewer';
  }[];
}

function TeamStats({ users }: TeamStatsProps) {
  return (
    <Card className="col-span-1">
      <CardHeader>
        <CardTitle>Organization Stats</CardTitle>
        <CardDescription>Overview of your organization</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          <div className="flex justify-between items-center">
            <span className="text-sm text-muted-foreground">Total Members</span>
            <span className="font-medium">{users.length}</span>
          </div>
          <div className="flex justify-between items-center">
            <span className="text-sm text-muted-foreground">Admins</span>
            <span className="font-medium">{users.filter((u) => u.role === 'Admin').length}</span>
          </div>
          <div className="flex justify-between items-center">
            <span className="text-sm text-muted-foreground">Members</span>
            <span className="font-medium">{users.filter((u) => u.role === 'Member').length}</span>
          </div>
          <div className="flex justify-between items-center">
            <span className="text-sm text-muted-foreground">Viewers</span>
            <span className="font-medium">{users.filter((u) => u.role === 'Viewer').length}</span>
          </div>
        </div>
      </CardContent>
      <CardFooter>
        <Button variant="outline" className="w-full">
          View Details
        </Button>
      </CardFooter>
    </Card>
  );
}

export default TeamStats;
