import React from 'react'
import EventItem from "./Eventitem"
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faExclamationCircle } from '@fortawesome/free-solid-svg-icons'

const EventsModal = () => {
  return (
    <button className="events-icon" href="#">
        <FontAwesomeIcon icon={faExclamationCircle}/>
        <table className="events">
            <tbody>
                <EventItem team="PHI" date="Mar 30" playername="playerName" updateMsg="Traded"/>
                <EventItem team="PHI" date="Apr 2" playername="playerName" updateMsg="Injured"/>
                <EventItem team="PHI" date="Apr 2" playername="playerName" updateMsg="Returned"/>
                <EventItem team="BKN" date="Apr 2" playername="playerName" updateMsg="Traded"/>
            </tbody>
        </table>
    </button>
  )
}

export default EventsModal