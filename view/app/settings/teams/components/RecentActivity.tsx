import React from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useGetRecentAuditLogsQuery } from '@/redux/services/audit';
import { formatDistanceToNow } from 'date-fns';
import { Loader2 } from 'lucide-react';
import { AuditAction, AuditLog } from '@/redux/types/audit';

const getActionColor = (action: AuditAction) => {
  switch (action) {
    case 'create':
      return 'bg-green-500';
    case 'update':
      return 'bg-blue-500';
    case 'delete':
      return 'bg-red-500';
    default:
      return 'bg-gray-500';
  }
};

const getActionMessage = (log: AuditLog) => {
  const username = log.user?.username || 'A user';
  const resource = log.resource_type
    .split('_')
    .map((word: string) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ');

  switch (log.action) {
    case 'create':
      return `${username} created a new ${resource.toLowerCase()}`;
    case 'update':
      return `${username} updated a ${resource.toLowerCase()}`;
    case 'delete':
      return `${username} deleted a ${resource.toLowerCase()}`;
    default:
      return `${username} accessed a ${resource.toLowerCase()}`;
  }
};

function RecentActivity() {
  const { data: auditLogs, isLoading, error } = useGetRecentAuditLogsQuery();
  return (
    <Card>
      <CardHeader>
        <CardTitle>Recent Activities</CardTitle>
        <CardDescription>Recent changes to your organization</CardDescription>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="flex items-center justify-center p-4">
            <Loader2 className="h-6 w-6 animate-spin" />
          </div>
        ) : error ? (
          <div className="p-4 text-red-600">Failed to load recent activities</div>
        ) : auditLogs && auditLogs.length > 0 ? (
          <div className="space-y-4">
            {auditLogs.map((log) => (
              <div key={log.id} className="flex items-start gap-4">
                <div className={`h-2 w-2 mt-2 rounded-full ${getActionColor(log.action)}`}></div>
                <div>
                  <p className="text-sm font-medium">{getActionMessage(log)}</p>
                  <p className="text-xs text-muted-foreground">
                    {formatDistanceToNow(new Date(log.created_at), { addSuffix: true })}
                  </p>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="text-center text-muted-foreground">No recent activities found</div>
        )}
      </CardContent>
    </Card>
  );
}

export default RecentActivity;
