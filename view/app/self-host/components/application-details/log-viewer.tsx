'use client';

import React from 'react';
import { Card, CardHeader, CardContent, CardDescription } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select';
import { Maximize2, Minimize2, Search, Download, RefreshCw, MoreHorizontal } from 'lucide-react';
import './logViewer.css';
import AceEditorComponent from '@/components/ui/ace-editor';
import useLogViewer, { LogViewerProps, LOG_LEVELS } from '../../hooks/use_log_viewer';
import { useTranslation } from '@/hooks/use-translation';

function LogViewer({
  logs = '',
  title = 'Log Viewer',
  description = 'Real-time installation logs',
  onRefresh,
  currentPage = 1,
  setCurrentPage
}: LogViewerProps) {
  const { t } = useTranslation();
  const {
    filteredLogs,
    handleRefresh,
    isRefreshing,
    handleEditorLoad,
    toggleFullscreen,
    isFullscreen,
    handleDownload,
    navigateSearch,
    currentSearchIndex,
    markers,
    searchTerm,
    setSearchTerm,
    selectedLevel,
    setSelectedLevel,
    timeRange,
    setTimeRange,
    autoScroll,
    setAutoScroll
  } = useLogViewer({
    logs,
    title,
    description,
    onRefresh,
    currentPage,
    setCurrentPage
  });

  return (
    <div
      className={`transition-all duration-300 ease-in-out ${isFullscreen ? 'fixed inset-0 z-50 bg-background' : 'relative w-full'}`}
    >
      <Card className={`h-full ${isFullscreen ? 'rounded-none' : ''}`}>
        <CardHeader className="flex flex-row items-center justify-between">
          <div>
            <h3 className="text-lg font-semibold">{t('selfHost.logViewer.title')}</h3>
            <CardDescription>{t('selfHost.logViewer.description')}</CardDescription>
          </div>
          <div className="flex items-center space-x-2">
            <Button variant="outline" onClick={() => setCurrentPage(currentPage + 1)}>
              <MoreHorizontal className="mr-2 h-4 w-4 text-muted-foreground" />
              {t('selfHost.logViewer.actions.fetchMore')}
            </Button>
            <Button variant="outline" size="icon" onClick={handleRefresh} disabled={isRefreshing}>
              <RefreshCw className={`h-4 w-4 ${isRefreshing ? 'animate-spin' : ''}`} />
            </Button>
            <Button variant="outline" size="icon" onClick={toggleFullscreen}>
              {isFullscreen ? <Minimize2 className="h-4 w-4" /> : <Maximize2 className="h-4 w-4" />}
            </Button>
          </div>
        </CardHeader>
        <CardContent className={`${isFullscreen ? 'h-[calc(100vh-120px)]' : ''}`}>
          <div className="h-full space-y-4">
            <div className="flex flex-row justify-between">
              <div className="flex items-center space-x-4">
                <div className="relative flex-1">
                  <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
                  <Input
                    placeholder={t('selfHost.logViewer.actions.search.placeholder')}
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    className="pl-8"
                  />
                </div>
                {markers.length > 0 && (
                  <div className="flex items-center space-x-2">
                    <Button variant="outline" size="sm" onClick={() => navigateSearch('prev')}>
                      {t('selfHost.logViewer.actions.search.previous')}
                    </Button>
                    <span className="text-sm">
                      {t('selfHost.logViewer.actions.search.results')
                        .replace('{current}', (currentSearchIndex + 1).toString())
                        .replace('{total}', markers.length.toString())}
                    </span>
                    <Button variant="outline" size="sm" onClick={() => navigateSearch('next')}>
                      {t('selfHost.logViewer.actions.search.next')}
                    </Button>
                  </div>
                )}
              </div>
              <div className="flex items-center space-x-4">
                <Select value={selectedLevel} onValueChange={setSelectedLevel}>
                  <SelectTrigger className="w-[180px]">
                    <SelectValue placeholder={t('selfHost.logViewer.actions.filter.level.label')} />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">
                      {t('selfHost.logViewer.actions.filter.level.all')}
                    </SelectItem>
                    {LOG_LEVELS?.map((level) => (
                      <SelectItem key={level.value} value={level.value} className={level.color}>
                        {level.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>

                <Select value={timeRange} onValueChange={setTimeRange}>
                  <SelectTrigger className="w-[180px]">
                    <SelectValue placeholder={t('selfHost.logViewer.actions.filter.time.label')} />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">
                      {t('selfHost.logViewer.actions.filter.time.all')}
                    </SelectItem>
                    <SelectItem value="5m">
                      {t('selfHost.logViewer.actions.filter.time.5m')}
                    </SelectItem>
                    <SelectItem value="15m">
                      {t('selfHost.logViewer.actions.filter.time.15m')}
                    </SelectItem>
                    <SelectItem value="1h">
                      {t('selfHost.logViewer.actions.filter.time.1h')}
                    </SelectItem>
                    <SelectItem value="24h">
                      {t('selfHost.logViewer.actions.filter.time.24h')}
                    </SelectItem>
                  </SelectContent>
                </Select>
                <div className="ml-auto flex items-center space-x-2">
                  <label className="flex items-center space-x-2">
                    <input
                      type="checkbox"
                      checked={autoScroll}
                      onChange={(e) => setAutoScroll(e.target.checked)}
                      className="form-checkbox"
                    />
                    <span>{t('selfHost.logViewer.actions.autoScroll')}</span>
                  </label>
                </div>
                <Button variant="outline" size="icon" onClick={handleDownload}>
                  <Download className="h-4 w-4" />
                </Button>
              </div>
            </div>
            <div className="h-[calc(100%-120px)]">
              <AceEditorComponent
                mode="sh"
                value={filteredLogs}
                onChange={() => {}}
                name="log-editor"
                readOnly={true}
                onLoad={handleEditorLoad}
                height={isFullscreen ? '100%' : '600px'}
                markers={markers}
              />
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

export default LogViewer;
