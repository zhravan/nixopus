import React$1 from 'react';
import { RectResult } from '@reactour/utils';

type StylesKeys = 'popover';
type StylesObj = {
    [key in StylesKeys]?: StyleFn;
};
type StyleFn = (props: {
    [key: string]: any;
}, state?: {
    [key: string]: any;
}) => React.CSSProperties;

declare const Popover: React$1.FC<PopoverProps>;

type PopoverProps = {
    sizes: RectResult;
    children?: React$1.ReactNode;
    position?: PositionType;
    padding?: number | number[];
    styles?: StylesObj;
    className?: string;
    refresher?: any;
};
type PositionType = Position | ((postionsProps: PositionProps, prevRect: RectResult) => Position);
type PositionProps = RectResult & {
    windowWidth: number;
    windowHeight: number;
};
type Position = 'top' | 'right' | 'bottom' | 'left' | 'center' | [number, number];

export { Popover, type StylesObj as PopoverStylesObj, type PositionType as Position, type PositionProps, Popover as default };
