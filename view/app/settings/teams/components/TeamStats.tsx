import React from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

interface TeamStatsProps {
  users: {
    id: string;
    name: string;
    role: 'admin' | 'member' | 'viewer' | 'owner';
  }[];
}

function TeamStats({ users }: TeamStatsProps) {
  return (
    <Card className="col-span-1">
      <CardHeader>
        <CardTitle>Team Stats</CardTitle>
        <CardDescription>Overview of your team</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          <div className="flex justify-between items-center">
            <span className="text-sm text-muted-foreground">Total Members</span>
            <span className="font-medium">{users.length}</span>
          </div>
          <div className="flex justify-between items-center">
            <span className="text-sm text-muted-foreground">Owners</span>
            <span className="font-medium">{users.filter((u) => u.role === 'admin').length}</span>
          </div>
          <div className="flex justify-between items-center">
            <span className="text-sm text-muted-foreground">Members</span>
            <span className="font-medium">{users.filter((u) => u.role === 'member').length}</span>
          </div>
          <div className="flex justify-between items-center">
            <span className="text-sm text-muted-foreground">Viewers</span>
            <span className="font-medium">{users.filter((u) => u.role === 'viewer').length}</span>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

export default TeamStats;
