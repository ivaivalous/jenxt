/*
<jenxt>
{
    "expose": "scm-triggers",
    "jsonResponse": true
}
</jenxt>
*/

import groovy.json.*

def result = [:]

Jenkins.instance.getAllItems().each { it ->
  if (!(it instanceof jenkins.triggers.SCMTriggerItem)) {
    return
  }

  def itTrigger = (jenkins.triggers.SCMTriggerItem)it
  def triggers = itTrigger.getSCMTrigger()

  triggers.each { t->
    def builder = new JsonBuilder()
    result[it.name] = builder {
      spec "${t.getSpec()}"
      ignorePostCommitHooks "${t.isIgnorePostCommitHooks()}"
    }
  }
}

return new JsonBuilder(result).toPrettyString()