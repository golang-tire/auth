import React from 'react';
import { Card, Col, Row } from 'antd';

const Dashboard = props => {
    return (
        <Row gutter={16}>
            <Col span={6}>
                <Card title="Users" bordered={false}>
                    Card content
                </Card>
            </Col>
            <Col span={6}>
                <Card title="Active Sessions" bordered={false}>
                    Card content
                </Card>
            </Col>
            <Col span={6}>
                <Card title="Alerts" bordered={false}>
                    Card content
                </Card>
            </Col>
            <Col span={6}>
                <Card title="New Users" bordered={false}>
                    Card content
                </Card>
            </Col>
        </Row>
    );
};

export default Dashboard;