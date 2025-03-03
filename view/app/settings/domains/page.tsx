'use client';
import DashboardPageHeader from '@/components/dashboard-page-header';
import { Button } from '@/components/ui/button';
import React from 'react';
import DomainsTable from './components/domainsTable';
import { ServerAndDomain } from '@/redux/types/domain';
import AddDomainDialog from './components/add-domain';

const mockData: { serverAndDomains: ServerAndDomain[] } = {
    serverAndDomains: [
        {
            id: '1',
            server_id: 'server-1',
            domain_id: 'domain-1',
            server: {
                id: 'server-1',
                name: 'Production Server',
                created_at: '2025-01-15T08:30:00Z',
                updated_at: '2025-02-20T14:45:00Z',
                deleted_at: null,
                is_primary: true
            },
            domains: [
                {
                    id: 'domain-1',
                    domain: 'example.com',
                    created_at: '2025-01-15T08:30:00Z',
                    updated_at: '2025-02-20T14:45:00Z',
                    deleted_at: null,
                    is_wildcard: false
                }
            ],
            created_at: '2025-01-15T08:30:00Z',
            updated_at: '2025-02-20T14:45:00Z',
            deleted_at: null
        }
    ]
};

const Page = () => {
    const { serverAndDomains } = mockData;
    const [addDomainDialogOpen, setAddDomainDialogOpen] = React.useState(false);

    return (
        <div className="container mx-auto py-6 space-y-8 max-w-4xl">
            <DashboardPageHeader
                label="Server and Domains"
                description="Configure your domains"
            />
            <div className="flex justify-between items-center mt-8">
                <h2 className="text-xl font-medium text-foreground">Domains</h2>
                <Button variant="default" onClick={() => setAddDomainDialogOpen(true)}>
                    {' '}
                    Add Domain
                </Button>
            </div>
            <DomainsTable serverAndDomains={serverAndDomains} />
            {addDomainDialogOpen && (
                <AddDomainDialog open={addDomainDialogOpen} setOpen={setAddDomainDialogOpen} />
            )}
        </div>
    );
};

export default Page;
