import React from 'react'

const EventItem = (props) => {
    const {team, date, playername, updateMsg} = props
  return (
    <tr><th><span className="team-abr bold">{team}</span></th><td>{playername}</td><td>{updateMsg}</td><td>{date}</td></tr>
  )
}

export default EventItem