import React, { useEffect, useState} from 'react';
import {Divider, Flex, Layout, Spin} from 'antd';
import ContainersList from "./ui/containersList/ContainersList.tsx";
import { Typography } from 'antd';
import {useFetching} from "./core/hooks/useFetching.tsx";
const {  Title } = Typography;
import { IContainer} from "./types/types.ts";
import { LoadingOutlined } from '@ant-design/icons';
import GetAll from "./core/API/getAll.ts";
import GetHistory from "./core/API/getHistory.ts";

const { Header, Content } = Layout;

const headerStyle: React.CSSProperties = {
    textAlign: 'center',
    color: '#fff',
    height: 128,
    paddingInline: 48,
    lineHeight: '64px',
    backgroundColor: '#fff',
};

const contentStyle: React.CSSProperties = {
    minHeight: 450,
    color: '#fff',
    backgroundColor: '#ffff',
};


const layoutStyle = {
    borderRadius: 8,
    overflow: 'hidden',

};

const App: React.FC = ({ }) => {
    const [containers, setContainers] = useState<IContainer[]>([]);
    const [history, setHistory] = useState<IContainer[]>([]);

    const [fetchContainers, isContainerLoading, _] = useFetching(async () => {
        const response = await GetAll();
        setContainers(response);
    });

    const [fetchHistory, isHistoryLoading] = useFetching(async () => {
        const response = await GetHistory();
        setHistory(response);
    });


    useEffect(() => {
        fetchContainers();
        fetchHistory();

    }, []);



    return(
        <Flex gap="middle" wrap>
            <Layout style={layoutStyle}>
                <Header style={headerStyle}>
                    <Title>Pinger</Title>
                </Header>
                <Content style={contentStyle}>
                    <>
                        <Title level={4}>Docker-контейнеры</Title>
                        <p style={{color: "black", textAlign: "left"}}>
                            Здесь отображаются контейнеры, которые пингуются раз в 5 секунд.
                            Можно изменить это значение в настройках
                        </p>
                        {isContainerLoading ? (
                            <Flex align="center" gap="middle">
                                <Spin indicator={<LoadingOutlined spin />} size="large" />
                            </Flex>
                        ) : <ContainersList data={containers} />}
                    </>

                    <>
                        <Divider orientation="left">История</Divider>
                        <p style={{color: "black", textAlign: "left"}}>
                            История всех запросов
                        </p>
                    </>
                    {isHistoryLoading ? (
                        <Flex align="center" gap="middle">
                            <Spin indicator={<LoadingOutlined spin />} size="large" />
                        </Flex>
                    ) : (
                        <ContainersList data={history} />
                    )}
                </Content>

            </Layout>

        </Flex>
    )
}


export default App;
