import { useEffect, useState } from 'react';
import Moment from 'react-moment';
import moment from 'moment';

export function VehicleList() {
  const [vehicles, setVehicles] = useState([])

  useEffect(()=> {
    fetch('/vehicles', { 'method' : 'get' })
      .then(response => response.json())
      .then(vehicles => setVehicles(vehicles))
  }, []);

  function vehicleComponents() {
    return vehicles.map((vehicle, i) => {
      return(<Vehicle key={vehicle.ID} {...vehicle} />);
    })
  }

  return(
    <div className='vehicleList'>
      {vehicleComponents()}
    </div>
  )
}

function Vehicle({RegistrationNumber, Manufacturer, Model, MotDue, VEDDue, MOTHistory}) {
  function expiredOrDue(timestamp) {
    if(moment(timestamp).isBefore(moment())){
      return 'Expired';
    } else {
      return 'Due';
    }
  }

  function latestMOT() {
    if(MOTHistory == null) {
      return null;
    }

    const sortedHistory = MOTHistory.sort((a, b) => {
      if(a.CompletedDate > b.CompletedDate) {
        return -1;
      } else if (a.CompletedDate < b.CompletedDate) {
        return 1;
      } else {
        return 0;
      }
    });

    return sortedHistory[0];
  }

  function advisoryCount() {
    const latest = latestMOT();

    if(latest != null) {
      const advisories = latest.RfrAndComments.filter(comment => comment.Type === 'MINOR')

      return advisories.length;
    }
  }

  return(
    <div className = 'vehicle'>
      <h2>{RegistrationNumber}</h2>
      <h3>{Manufacturer}</h3>
      <h3>{Model}</h3>
      <dl>
        <dt>VED Status</dt>
        <dd>{expiredOrDue(VEDDue)} <Moment format='DD/MM/YYYY'>{VEDDue}</Moment></dd>
        <dt>MOT Status</dt>
        <dd>{expiredOrDue(MotDue)} <Moment format='DD/MM/YYYY'>{MotDue}</Moment></dd>
        <dt>Advisory Count</dt>
        <dd>{advisoryCount()}</dd>
      </dl>
    </div>
  );
}
