export interface IContainer{
    addr: string;
    container_status: string;
    p_duration: number;
    pinged_at: string;
}
export interface FetchResponse {
    total_count: number;
    total_pages: number;
    page: number;
    size: number;
    has_more: boolean;
    containers: IContainer[];
}
