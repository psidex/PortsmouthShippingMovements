export function fromToSentence(movement) {
  if (movement.type === 0) {
    return `Moving from ${movement.from.name} to ${movement.to.name}`;
  }
  return '';
}

export function fromToAbbreviation(movement) {
  if (movement.type === 0) {
    return `${movement.from.abbreviation} to ${movement.to.abbreviation}`;
  }
  return '';
}

export async function getMovements() {
  let r;

  if (!process.env.NODE_ENV || process.env.NODE_ENV === 'development') {
    r = await fetch('./exampledata.json');
  } else {
    r = await fetch('/api/movements');
  }

  const movements = await r.json();

  if (movements.today === null) {
    movements.today = [];
  }
  if (movements.tomorrow === null) {
    movements.tomorrow = [];
  }

  return movements;
}
