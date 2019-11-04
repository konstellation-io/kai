import React, {useState} from 'react';
import PropTypes from 'prop-types';
import Button, {BUTTON_TYPES, BUTTON_ALIGN} from '../Button/Button';
import { ICON } from '../.././icons';
import styles from './Settings.module.scss';

const BUTTON_HEIGHT = 40;

const buttonStyle = {
  paddingLeft: '20%'
};

function doLogout() {
  console.log('LOGGED OUT');
}

function Settings({
  label = 'Button',
} = {}) {
  const [opened, setOpened] = useState(false);

  const buttons = [
    <Button
      label={'LOGOUT'}
      type={BUTTON_TYPES.GREY}
      onClick={doLogout}
      icon={ICON.LOGOUT}
      align={BUTTON_ALIGN.LEFT}
      style={buttonStyle}
      key={'buttonLogout'}
    />
  ];
  const nButtons = buttons.length;
  const optionsHeight = nButtons * BUTTON_HEIGHT;

  return (
    <div
      className={styles.container}
      onMouseEnter={ () => setOpened(true) }
      onMouseLeave={ () => setOpened(false) }
      data-testid='settingsContainer'
    >
      <div className={styles.label}>
        { label }
        <div className={styles.arrow} />
      </div>
      <div
        className={styles.options}
        style={{ maxHeight: opened ? optionsHeight : 0 }}
        data-testid='settingsContent'
      >
        { buttons }
      </div>
    </div>
  );
}

Settings.propTypes = {
  label: PropTypes.string,
};

export default Settings;
