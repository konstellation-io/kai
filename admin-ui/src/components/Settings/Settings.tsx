import { get } from 'lodash';
import React, { useState, useEffect } from 'react';
import { useHistory } from 'react-router';
import useEndpoint from '../../hooks/useEndpoint';
import { useApolloClient } from '@apollo/react-hooks';

import { ENDPOINT } from '../../constants/application';

import Button, { BUTTON_TYPES, BUTTON_ALIGN } from '../Button/Button';
import SettingsIcon from '@material-ui/icons/VerifiedUser';
import AuditIcon from '@material-ui/icons/SupervisorAccount';
import LogoutIcon from '@material-ui/icons/ExitToApp';

import styles from './Settings.module.scss';
import ROUTE from '../../constants/routes';

const BUTTON_HEIGHT = 40;
const buttonStyle = {
  paddingLeft: '20%'
};

type Props = {
  label: string;
};

function Settings({ label }: Props) {
  const client = useApolloClient();
  const history = useHistory();
  const [opened, setOpened] = useState(false);
  const [logoutResponse, doLogout] = useEndpoint({
    endpoint: ENDPOINT.LOGOUT,
    method: 'POST'
  });

  useEffect(() => {
    if (logoutResponse.complete) {
      if (get(logoutResponse, 'status') === 200) {
        client.writeData({ data: { loggedIn: false } });
        history.push(ROUTE.LOGIN);
      } else {
        console.error(`Error sending logout request.`);
      }
    }
  }, [logoutResponse, history, client]);

  function getBaseProps(label: string) {
    return {
      label: label.toUpperCase(),
      key: `button${label}`,
      type: BUTTON_TYPES.GREY,
      align: BUTTON_ALIGN.LEFT,
      style: buttonStyle
    };
  }

  const buttons = [
    <Button
      {...getBaseProps('Settings')}
      to={ROUTE.SETTINGS}
      Icon={SettingsIcon}
    />,
    <Button {...getBaseProps('Audit')} to={ROUTE.AUDIT} Icon={AuditIcon} />,
    <Button
      {...getBaseProps('Logout')}
      onClick={() => doLogout()}
      Icon={LogoutIcon}
    />
  ];
  const nButtons = buttons.length;
  const optionsHeight = nButtons * BUTTON_HEIGHT;

  return (
    <div
      className={styles.container}
      onMouseEnter={() => setOpened(true)}
      onMouseLeave={() => setOpened(false)}
      data-testid="settingsContainer"
    >
      <div className={styles.label}>
        {label}
        <div className={styles.arrow} />
      </div>
      <div
        className={styles.options}
        style={{ maxHeight: opened ? optionsHeight : 0 }}
        data-testid="settingsContent"
      >
        {buttons}
      </div>
    </div>
  );
}

export default Settings;
