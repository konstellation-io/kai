import React, { useState } from 'react';
import moment from 'moment';

import IconInfo from '@material-ui/icons/Info';
import IconExpand from '@material-ui/icons/ArrowDownward';

import styles from './LogsList.module.scss';
import cx from 'classnames';
import { GetServerLogs_logs_items } from '../../../../../../../graphql/queries/types/GetServerLogs';

function LogItem({ date, nodeName, message, level }: GetServerLogs_logs_items) {
  const [opened, setOpened] = useState<boolean>(false);

  const dateFormatted = moment(date).format('YYYY-MM-DD');
  const hourFormatted = moment(date).format('HH:mm:ss');

  const LevelIcon = IconInfo;

  function toggleOpenStatus() {
    setOpened(!opened);
  }

  return (
    <div
      className={cx(styles.container, {
        [styles.opened]: opened
      })}
    >
      <div className={styles.row1}>
        <div className={styles.icon}>
          <LevelIcon className="icon-small" />
        </div>
        <div className={styles.date}>{dateFormatted}</div>
        <div className={styles.hour}>{hourFormatted}</div>
        {/* Uncomment this when adding workflows and add workflows and toggle logic */}
        {/* <div className={styles.nodeName}>{nodeName}</div> */}
        <div className={styles.message}>{message}</div>
        <div className={styles.expand} onClick={toggleOpenStatus}>
          <IconExpand className="icon-regular" />
        </div>
      </div>
      {opened && <div className={styles.messageComplete}>{message}</div>}
    </div>
  );
}

export default LogItem;