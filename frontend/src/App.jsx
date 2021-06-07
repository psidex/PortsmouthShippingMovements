import './css/App.scss';
import 'react-loader-spinner/dist/loader/css/react-spinner-loader.css';

import React, { useEffect, useState } from 'react';
import Loader from 'react-loader-spinner';

import ShipMovements from './ShipMovements';
import { getMovements } from './api/movements';

export default function App() {
  const [movements, setMovements] = useState({});
  const [loaded, setLoaded] = useState(false);

  useEffect(async () => {
    setMovements(await getMovements());
    setLoaded(true);
  }, []);

  if (!loaded) {
    return (
      <Loader
        className="loader"
        type="Grid"
        color="#00BFFF"
        height={100}
        width={100}
      />
    );
  }

  return (
    <div id="app">
      <h1>Portsmouth Shipping Movements</h1>
      <ShipMovements id="todayMovements" title="Today" movements={movements.today} />
      <ShipMovements id="tomorrowMovements" title="Tomorrow" movements={movements.tomorrow} />
    </div>
  );
}
