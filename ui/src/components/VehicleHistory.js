import { useEffect, useState } from 'react';
import { useParams } from "react-router-dom";
import Moment from 'react-moment';

export default function VehicleHistory() {
  const { id } = useParams();
  const [vehicle, setVehicle] = useState(null);

  useEffect(()=> {
    fetch('/vehicles', { 'method' : 'get' })
      .then(response => response.json())
      .then(vehicles => findVehicle(vehicles))
      .then(vehicle => setVehicle(vehicle))
  }, [id]);

  function findVehicle(vehicles) {
    return vehicles.filter(vehicle => vehicle.ID === id)[0];
  }

  function showVehicle() {
    if(vehicle === null) {
      return(<h2>Loading</h2>)
    } else {
      console.log(vehicle)
      return(<Vehicle {...vehicle} />)
    }
  }

  return(showVehicle())
}

function Vehicle({Manufacturer, Model, RegistrationNumber, MOTHistory }) {
  function MOTs() {
    return MOTHistory.map(mot => <MOTTest key={mot.TestNumber} {...mot}/>)
  }

  return(
    <div>
      <h2>{Manufacturer} {Model} - {RegistrationNumber}</h2>
      {MOTs()}
    </div>
  );
}

function MOTTest({ Passed, CompletedDate }) {
  return(
    <div>
      <p>
        {Passed ? 'Passed' : 'Failed'}
        <br/>
        <Moment>{CompletedDate}</Moment>
      </p>
    </div>
  )
}
