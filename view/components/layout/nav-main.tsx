'use client';

import { ChevronRight, type LucideIcon } from 'lucide-react';
import { usePathname, useRouter } from 'next/navigation';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible';
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem
} from '@/components/ui/sidebar';
import Link from 'next/link';
import { useCollapsibleState } from '@/hooks/use-collapsible-state';

interface NavItem {
  title: string;
  url: string;
  icon?: LucideIcon;
  isActive?: boolean;
  items?: { title: string; url: string }[];
}

interface NavMainProps {
  items: NavItem[];
  onItemClick?: (url: string) => void;
}

export function NavMain({ items, onItemClick }: NavMainProps) {
  const pathname = usePathname();
  const router = useRouter();
  const { isItemCollapsed, toggleItem } = useCollapsibleState();

  const handleClick = (url: string) => {
    onItemClick?.(url);
    router.push(url);
  };

  return (
    <SidebarGroup>
      <SidebarGroupLabel>Platform</SidebarGroupLabel>
      <SidebarMenu>
        {items.map((item) => (
          <Collapsible
            key={item.title}
            asChild
            open={!isItemCollapsed(item.title)}
            onOpenChange={() => toggleItem(item.title)}
            className="group/collapsible"
          >
            <SidebarMenuItem>
              <CollapsibleTrigger asChild>
                <SidebarMenuButton tooltip={item.title} onClick={() => handleClick(item.url)}>
                  {item.icon && <item.icon />}
                  <span>{item.title}</span>
                  {(item.items?.length || 0) > 0 && (
                    <ChevronRight className="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90" />
                  )}
                </SidebarMenuButton>
              </CollapsibleTrigger>
              {(item.items?.length || 0) > 0 && (
                <CollapsibleContent>
                  <SidebarMenuSub>
                    {item.items?.map((subItem) => (
                      <SidebarMenuSubItem key={subItem.title}>
                        <SidebarMenuSubButton asChild>
                          <Link href={subItem.url}>
                            <span>{subItem.title}</span>
                          </Link>
                        </SidebarMenuSubButton>
                      </SidebarMenuSubItem>
                    ))}
                  </SidebarMenuSub>
                </CollapsibleContent>
              )}
            </SidebarMenuItem>
          </Collapsible>
        ))}
      </SidebarMenu>
    </SidebarGroup>
  );
}
