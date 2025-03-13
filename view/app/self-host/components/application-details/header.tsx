'use client';
import React from 'react';
import { Button } from '@/components/ui/button';
import { ExternalLink, MoreVertical, RotateCcw } from 'lucide-react';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Tooltip, TooltipProvider, TooltipTrigger, TooltipContent } from '@/components/ui/tooltip';
import { Application } from '@/redux/types/applications';
import { DeleteDialog } from '@/components/delete-dialog';

const ApplicationDetailsHeader = ({
    application
}: {
    application?: Application
}) => {
    return (
        <div className="flex flex-col items-start justify-between gap-4 sm:flex-row sm:items-center">
            <div className="flex items-start">
                <div className="mr-2">
                    <div className="flex items-center gap-2">
                        <h1 className="text-2xl font-bold capitalize sm:text-3xl">
                            {application?.name}
                        </h1>
                        <Button
                            variant="ghost"
                            size="icon"
                            onClick={() =>
                                window.open('https://' + application?.domain?.name, '_blank')
                            }
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
                    jobName={application?.name || ''}
                    onDelete={() => {}}
                    showButton={false}
                    isDeleting={false}
                />
                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                        <Button
                            variant="outline"
                            size="icon"

                        >
                            <MoreVertical className="h-4 w-4" />
                        </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                        <DropdownMenuItem
                        >
                            Force Deploy Without Cache
                        </DropdownMenuItem>
                        <DropdownMenuItem
                        >
                            Force Deploy
                        </DropdownMenuItem>
                    </DropdownMenuContent>
                </DropdownMenu>
            </div>
        </div>
    );
};

export default ApplicationDetailsHeader;