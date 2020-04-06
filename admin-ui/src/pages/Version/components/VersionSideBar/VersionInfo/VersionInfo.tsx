import React from 'react';
import { GetVersionConfStatus_versions } from '../../../../../graphql/queries/types/GetVersionConfStatus';
import StatusCircle from '../../../../../components/StatusCircle/StatusCircle';
import cx from 'classnames';
import styles from './VersionInfo.module.scss';

type VersionListItemProps = {
  version: GetVersionConfStatus_versions;
};

function VersionInfo({ version }: VersionListItemProps) {
  return (
    <div className={styles.wrapper}>
      <div className={styles.name}>
        <div>{version.name}</div>
      </div>
      <div className={styles.labelContainer}>
        <StatusCircle status={version.status} />
        <div className={cx(styles.label, styles[version.status])}>
          {version.status}
        </div>
      </div>

      <div>
        <div className={styles.desc}>{version.description}</div>
      </div>
    </div>
  );
}

export default VersionInfo;
