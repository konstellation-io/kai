import React from 'react';
import Button from '../../Button/Button';
import IconClose from '@material-ui/icons/Close';
import { DoubleSelectData } from './DoubleSelect';
import styles from './DoubleSelect.module.scss';

type SelectionProps = {
  group: string;
  element: string;
  onDeselect: (group: string, element: string) => void;
};

function Selection({ group, element, onDeselect }: SelectionProps) {
  const label = `${group}: ${element}`;
  return (
    <div className={styles.selection} title={label}>
      <div className={styles.label}>{label}</div>
      <div
        className={styles.selectionIcon}
        onClick={() => onDeselect(group, element)}
      >
        <IconClose className="icon-small" />
      </div>
    </div>
  );
}

type Props = {
  selections?: DoubleSelectData;
  onDeselect: (group: string, element: string) => void;
  onClear: () => void;
};

function Selections({ selections = {}, onClear, onDeselect }: Props) {
  const selectionsArray = Object.entries(selections)
    .map(([group, elements]) => elements.map(element => ({ group, element })))
    .flat();
  const selectionNodes = selectionsArray.map(({ group, element }) => (
    <Selection
      key={`${group}${element}`}
      group={group}
      element={element}
      onDeselect={onDeselect}
    />
  ));

  return (
    <div className={styles.selectionsContainer}>
      {selectionNodes}
      <div className={styles.clear} onClick={onClear}>
        <Button label="CLEAR" height={24} />
      </div>
    </div>
  );
}

export default Selections;
