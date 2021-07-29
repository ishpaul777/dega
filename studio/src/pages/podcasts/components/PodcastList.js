import React from 'react';
import { Popconfirm, Button, Table, Space } from 'antd';

import { useDispatch } from 'react-redux';
import { deletePodcast } from '../../../actions/podcasts';
import { Link } from 'react-router-dom';

function PodcastList({ actions, data, filters, setFilters, fetchPodcasts }) {
  const dispatch = useDispatch();
  const columns = [
    {
      title: 'Title',
      dataIndex: 'title',
      key: 'title',
      render: (_, record) => {
        return (
          <Link
            className="ant-dropdown-link"
            style={{
              marginRight: 8,
            }}
            to={`/podcasts/${record.id}/edit`}
          >
            {record.title}
          </Link>
        );
      },
    },
    { title: 'Slug', dataIndex: 'slug', key: 'slug' },
    { title: 'Season', dataIndex: 'season', key: 'season' },
    { title: 'Podcast', dataIndex: 'podcast', key: 'podcast' },
    {
      title: 'Action',
      dataIndex: 'operation',
      render: (_, record) => {
        return (
          <Popconfirm
            title="Are you sure you want to delete this?"
            onConfirm={() => dispatch(deletePodcast(record.id)).then(() => fetchPodcasts())}
          >
            <Link to="" className="ant-dropdown-link">
              <Button
                disabled={!(actions.includes('admin') || actions.includes('delete'))}
                type="danger"
              >
                Delete
              </Button>
            </Link>
          </Popconfirm>
        );
      },
    },
  ];

  return (
    <Space direction={'vertical'}>
      <Table
        bordered
        columns={columns}
        dataSource={data.podcasts}
        loading={data.loading}
        rowKey={'id'}
        pagination={{
          total: data.total,
          current: filters.page,
          pageSize: filters.limit,
          onChange: (pageNumber, pageSize) =>
            setFilters({ ...filters, page: pageNumber, limit: pageSize }),
        }}
      />
    </Space>
  );
}

export default PodcastList;
