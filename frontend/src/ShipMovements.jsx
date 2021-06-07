import './css/Accordion.scss';

import React from 'react';
import {
  Accordion,
  AccordionItem,
  AccordionItemButton,
  AccordionItemHeading,
  AccordionItemPanel,
} from 'react-accessible-accordion';

import { fromToAbbreviation, fromToSentence } from './api/movements';

export default function ShipMovements(props) {
  const { title, movements } = props;
  const items = [];

  for (const m of movements) {
    items.push(
      <AccordionItem>
        <AccordionItemHeading>
          <AccordionItemButton>
            <p>{m.time}</p>
            <p>{m.name}</p>
            <p>{fromToAbbreviation(m)}</p>
          </AccordionItemButton>
        </AccordionItemHeading>
        <AccordionItemPanel>
          <img src={m.imageUrl} alt={`The ship ${m.name}`} />
          <p>
            {fromToSentence(m)}
          </p>
        </AccordionItemPanel>
      </AccordionItem>,
    );
  }

  return (
    <div className="accordionDiv">
      <h2>{title}</h2>
      <Accordion allowMultipleExpanded allowZeroExpanded>
        {items}
      </Accordion>
    </div>
  );
}
