import React from 'react';
import { ServerAndDomain } from '@/redux/types/domain';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from '@/components/ui/table';
import { DomainTypeTag } from './domain-type-tag';
import { DomainActions } from './domain-actions';
import { ServerTag } from './server-tag';

interface DomainsTableProps {
  serverAndDomains: ServerAndDomain[];
}

function DomainsTable({ serverAndDomains }: DomainsTableProps) {
  return (
    <div className="bg-secondary shadow-md rounded-lg overflow-hidden p-4">
      <Table>
        <TableHeader className="bg-muted/50">
          <TableRow>
            <TableHead className="text-xs font-medium text-muted-foreground">Domain</TableHead>
            <TableHead className="text-xs font-medium text-muted-foreground">Type</TableHead>
            <TableHead className="text-xs font-medium text-muted-foreground">Server</TableHead>
            <TableHead className="text-right text-xs font-medium text-muted-foreground">
              Actions
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody className="bg-secondary divide-y divide-border">
          {serverAndDomains.flatMap((item) =>
            item.domains.map((domain) => (
              <TableRow key={domain.id}>
                <TableCell className="whitespace-nowrap font-medium text-foreground">
                  {domain.domain}
                </TableCell>
                <TableCell className="whitespace-nowrap text-muted-foreground">
                  <DomainTypeTag isWildcard={domain.is_wildcard} />
                </TableCell>
                <TableCell className="whitespace-nowrap text-muted-foreground">
                  <ServerTag name={item.server.name} isPrimary={item.server.is_primary} />
                </TableCell>
                <TableCell className="whitespace-nowrap text-right">
                  <DomainActions domainId={domain.id} />
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  );
}

export default DomainsTable;
