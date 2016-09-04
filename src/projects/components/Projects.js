import React, { Component } from 'react'

import ProjectItem from './ProjectItem'

class Projects extends Component {
  render() {
    return (
      <div className="projects-list">
        { this.props.projects.map(project =>
          <ProjectItem
            key={ project.id }
            project={ project }
            displayedPeriod={ this.props.displayedPeriod }
            toggleArchive={ () => this.props.toggleArchiveProject(project) }
            remove={ () => this.props.removeProject(project) }
            update={ data => this.props.updateProject(project, data) }
            updateMoodByWeek={ (weekNumber, data) => this.props.updateMoodByWeek(project, weekNumber, data) }
          />
        ) }
      </div>
    )
  }
}

export default Projects
