import React from 'react';
import { Space, Table, Tag} from 'antd';
import type { TableProps } from 'antd';
import {IContainer} from "../../types/types.ts";
import {format} from "date-fns";





type ContainerStatus = "created" | "running" | "paused" | "restarting" | "exited" | "dead";

type StatusColor = Record<ContainerStatus, string> & { [key: string]: string };

const statusColors: StatusColor = {
    created: "blue",
    running: "green",
    paused: "yellow",
    restarting: "orange",
    exited: "red",
    dead: "darkred",
    default: "gray",
};


const columns: TableProps<IContainer>['columns'] = [
    {
        title: "Адрес",
        dataIndex: 'addr',
        key: 'addr',
    },
    {
        title: "Статус",
        dataIndex: 'container_status',
        key: 'container_status',
        render: (_, record) => (
            <Tag color={statusColors[record.container_status as ContainerStatus] || statusColors.default}>
                {record.container_status}
            </Tag>
        ),
    },
    {
        title: "Время пинга (µs)",
        dataIndex: 'p_duration',
        key: 'p_duration',
        render: (value: number) => `${(value / 1000).toFixed(3)}`,
    },
    {
        title: "Последний пинг",
        dataIndex: 'pinged_at',
        key: 'pinged_at',
        render: (value: string) => format(new Date(value), "dd.MM.yyyy HH:mm:ss")
    },
];

interface ContainersListProps {
    data: IContainer[];
}

const ContainersList: React.FC<ContainersListProps> = ({ data }) => {
    return (
        <>
            <Space direction="vertical" style={{ marginBottom: '20px' }} size="middle" />
            <Table<IContainer>
                pagination={{
                    pageSize: 10,
                    position: ["bottomCenter"], // Расположение пагинации по центру снизу
                }}
                columns={columns}
                dataSource={data} />
        </>
    )
}



export default ContainersList;
