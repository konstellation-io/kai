import { get } from 'lodash';

import React from 'react';
import { formatDate } from '../../../../utils/format';

import CalendarIcon from '@material-ui/icons/Today';
import TimeIcon from '@material-ui/icons/AccessTime';

import cx from 'classnames';
import styles from './SidebarTitle.module.scss';

type DateInfoProps = {
  label: string;
  Icon: any;
  date?: string;
};
function DateInfo({ label, Icon, date }: DateInfoProps) {
  return (
    <>
      <div className={styles.dateTitle}>
        <Icon className="icon-small-regular" />
        <div>{label}</div>
      </div>
      <div className={styles.dateValue}>{date}</div>
    </>
  );
}

type Props = {
  version?: {
    title: string;
    status: string;
    name: string;
    created: string;
    activated: string;
  };
};

function SidebarTitle({ version }: Props) {
  if (!version) {
    return <div>No active version</div>; // TODO add final design
  }

  return (
    <div className={styles.container}>
      <div className={styles.title}>{get(version, 'title')}</div>
      <div className={styles.version}>
        <div
          className={cx(styles.versionStatus, {
            // @ts-ignore
            [styles[version.status]]: get(version, 'status')
          })}
        />
        <p>{`VERSION ${get(version, 'name')}`}</p>
      </div>
      <div className={styles.dates}>
        <DateInfo
          label={'CREATED'}
          date={version && formatDate(new Date(version.created))}
          Icon={CalendarIcon}
        />
        <DateInfo
          label={'ACTIVATED'}
          date={version && formatDate(new Date(version.activated))}
          Icon={TimeIcon}
        />
      </div>
    </div>
  );
}

export default SidebarTitle;
