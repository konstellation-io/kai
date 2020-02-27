import React, { useState, FunctionComponent } from 'react';
import {
  GetVersionConfStatus_runtime,
  GetVersionConfStatus_versions
} from '../../../../../graphql/queries/types/GetVersionConfStatus';
import { VersionStatus } from '../../../../../graphql/types/globalTypes';
import styles from './VersionActions.module.scss';

// Icons
import StartIcon from '@material-ui/icons/PlayArrowOutlined';
import StopIcon from '@material-ui/icons/Stop';
import PublishIcon from '@material-ui/icons/Publish';
import UnpublishIcon from '@material-ui/icons/Block';
import useVersionAction from '../../../utils/hooks';
import ConfirmationModal from '../../../../../components/ConfirmationModal/ConfirmationModal';
import Button, { BUTTON_TYPES } from '../../../../../components/Button/Button';
import { SvgIconProps } from '@material-ui/core/SvgIcon';

type VersionActionsProps = {
  runtime: GetVersionConfStatus_runtime;
  version: GetVersionConfStatus_versions;
};

type ActionProps = {
  Icon: FunctionComponent<SvgIconProps>;
  label: string;
  action?: Function;
  disabled?: boolean;
  primary?: boolean;
};

function VersionActions({ runtime, version }: VersionActionsProps) {
  const [showActionConfirmation, setShowPublishModal] = useState(false);
  const {
    publishVersion,
    startVersion,
    stopVersion,
    unpublishVersion,
    getMutationVars
  } = useVersionAction(runtime.id);

  function onPublishVersion(comment: string) {
    closePublishModal();
    publishVersion(getMutationVars(version.id, comment));
  }

  function onUnpublishVersion() {
    unpublishVersion(getMutationVars(version.id));
  }

  function onStartVersion() {
    startVersion(getMutationVars(version.id));
  }

  function onStopVersion() {
    stopVersion(getMutationVars(version.id));
  }

  function openPublishModal() {
    setShowPublishModal(true);
  }

  function closePublishModal() {
    setShowPublishModal(false);
  }

  const buttons: ActionProps[] = [];

  switch (version.status) {
    case VersionStatus.STOPPED:
      buttons[0] = {
        Icon: StartIcon,
        label: 'START',
        action: onStartVersion,
        primary: version.configurationCompleted,
        disabled: !version.configurationCompleted
      };
      buttons[1] = {
        Icon: PublishIcon,
        label: 'PUBLISH',
        disabled: true
      };
      break;
    case VersionStatus.STARTED:
      buttons[0] = {
        Icon: StopIcon,
        label: 'STOP',
        action: onStopVersion
      };
      buttons[1] = {
        Icon: PublishIcon,
        label: 'PUBLISH',
        action: openPublishModal,
        primary: true
      };
      break;
    case VersionStatus.PUBLISHED:
      buttons[0] = {
        Icon: StopIcon,
        label: 'STOP',
        disabled: true
      };
      buttons[1] = {
        Icon: UnpublishIcon,
        label: 'UNPUBLISH',
        action: onUnpublishVersion,
        primary: true
      };
      break;
  }

  return (
    <div className={styles.wrapper}>
      {buttons.map(btn => (
        <Button
          key={btn.label}
          onClick={() => btn.action && btn.action()}
          label={btn.label}
          Icon={btn.Icon}
          type={BUTTON_TYPES.DEFAULT}
          disabled={btn.disabled}
          primary={btn.primary}
          style={{ flexGrow: 1 }}
        />
      ))}
      {showActionConfirmation && (
        <ConfirmationModal
          title="YOU ARE ABOUT TO PUBLISH A VERSION"
          message=""
          onAction={onPublishVersion}
          onClose={closePublishModal}
        />
      )}
    </div>
  );
}

export default VersionActions;
