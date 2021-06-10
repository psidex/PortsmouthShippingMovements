import '../css/Accordion.scss';

import React from 'react';
import {
  Accordion,
  AccordionItem,
  AccordionItemButton,
  AccordionItemHeading,
  AccordionItemPanel,
} from 'react-accessible-accordion';

import { fromToAbbreviation, fromToSentence } from '../api/movements';

export default function ShipMovements(props) {
  const { id, title, movements } = props;
  const items = [];
  const startExpandedIds = [];

  for (const m of movements) {
    // AccordionItem UUID must not have whitespace.
    const uniqueId = `${m.time}-${m.name.replace(/\s/g, '')}`;

    const isNavy = m.name.includes('HMS');
    if (isNavy) {
      startExpandedIds.push(uniqueId);
    }

    items.push(
      <AccordionItem key={uniqueId} uuid={uniqueId}>
        <AccordionItemHeading>
          <AccordionItemButton>
            <p navy={isNavy ? 'true' : 'false'}>
              {m.time}
              {' - '}
              {m.name}
            </p>
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
    <div id={id} className="accordionDiv">
      <h2>{title}</h2>
      <Accordion allowMultipleExpanded allowZeroExpanded preExpanded={startExpandedIds}>
        {items}
      </Accordion>
    </div>
  );
}
