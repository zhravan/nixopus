import { Container } from '@/redux/services/container/containerApi';

export interface ContainerGroup {
  applicationId: string;
  applicationName: string;
  containers: Container[];
}

export interface GroupedContainers {
  groups: ContainerGroup[];
  ungrouped: Container[];
}

export function groupContainersByApplication(containers: Container[]): GroupedContainers {
  const groupsMap = new Map<string, ContainerGroup>();
  const ungrouped: Container[] = [];

  for (const container of containers) {
    const applicationId = container.labels?.['com.application.id'];
    const applicationName = container.labels?.['com.application.name'] || 'Unknown Application';

    if (applicationId) {
      if (!groupsMap.has(applicationId)) {
        groupsMap.set(applicationId, {
          applicationId,
          applicationName,
          containers: []
        });
      }
      groupsMap.get(applicationId)!.containers.push(container);
    } else {
      ungrouped.push(container);
    }
  }

  const groups = Array.from(groupsMap.values()).sort((a, b) =>
    a.applicationName.localeCompare(b.applicationName)
  );

  return { groups, ungrouped };
}
