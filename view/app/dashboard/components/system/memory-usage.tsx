'use client';

import React from 'react';
import { BarChart } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { SystemStatsType } from '@/redux/types/monitor';
import { Skeleton } from '@/components/ui/skeleton';
import { useTranslation } from '@/hooks/use-translation';

interface MemoryUsageCardProps {
  systemStats: SystemStatsType;
}

const MemoryUsageCard: React.FC<MemoryUsageCardProps> = ({ systemStats }) => {
  const { t } = useTranslation();
  const { memory } = systemStats;

  return (
    <Card className="overflow-hidden">
      <CardHeader className="pb-2">
        <CardTitle className="text-xs sm:text-sm font-medium flex items-center">
          <BarChart className="h-3 w-3 sm:h-4 sm:w-4 mr-1 sm:mr-2" />
          {t('dashboard.memory.title')}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-2 sm:space-y-3">
          <div className="w-full h-2 bg-gray-200 rounded-full">
            <div
              className={`h-2 rounded-full ${memory.percentage > 80 ? 'bg-red-500' : 'bg-green-500'}`}
              style={{ width: `${memory.percentage}%` }}
            />
          </div>
          <div className="flex justify-between text-xs text-muted-foreground">
            <span>{t('dashboard.memory.used').replace('{value}', memory.used.toString())}</span>
            <span>
              {t('dashboard.memory.percentage').replace('{value}', memory.percentage.toString())}
            </span>
            <span>{t('dashboard.memory.total').replace('{value}', memory.total.toString())}</span>
          </div>
          <div className="text-xs text-muted-foreground mt-1 sm:mt-2 line-clamp-2 sm:line-clamp-none">
            {memory.rawInfo}
          </div>
        </div>
      </CardContent>
    </Card>
  );
};

export default MemoryUsageCard;

export function MemoryUsageCardSkeleton() {
  const { t } = useTranslation();

  return (
    <Card className="overflow-hidden">
      <CardHeader className="pb-2">
        <CardTitle className="text-xs sm:text-sm font-medium flex items-center">
          <BarChart className="h-3 w-3 sm:h-4 sm:w-4 mr-1 sm:mr-2" />
          {t('dashboard.memory.title')}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-2 sm:space-y-3">
          <Skeleton className="w-full h-2 rounded-full" />

          <div className="flex justify-between text-xs">
            <Skeleton className="h-4 w-20" />
            <Skeleton className="h-4 w-10" />
            <Skeleton className="h-4 w-20" />
          </div>

          <Skeleton className="h-4 w-full mt-1 sm:mt-2" />
          <Skeleton className="h-4 w-2/3" />
        </div>
      </CardContent>
    </Card>
  );
}
