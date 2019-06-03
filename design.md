CI/CD tools in which you can write your own pipelines in golang code (idea from Gaia where you can write it in any code of your choosing)
The pro of this is you are not subjected to anyones shitty languages for writing your pipelines (I'm looking at you jenkins)
The con of this is that you have to build your own libraries for all the different actions you might want to take


## MVP Goals
* Master/Minion Arch: One master only, Potentially many minions
* Master will round robin jobs to minions
* Pipeline code will only be able to be imported from github
* Use Hashicorp's plugin framework to allow pipelines to be compiled and ran
* The ability to trigger a recompile of the repo
* Built using all GRPC
* Masters on first creation of pipeline will grab from github and compile
* It will run under the user of the program and store the compiled code in a specific directory
* Maybe use boltdb so that the database is self contained?
* Pipelines can be created from different github branches
* Workers will connect with a shared secret for now
* workers will have a keep alive thread. The master will check the last checkin time of a minion before
    it round robins a request to it
* workers can be "suspended" which means jobs will no longer be sent to them
* The worker list and data is kept in memory by the master
* How do workers communicate they are not in a healthy enough state to continue operation?


## Far Far Future Goals
* Master/Minion Arch with Masters able to be distributed by paxos
* Authenication and ACLs can build upon github groups and oauth
* Auto recompilation of github repo


## Questions
How do users use the filesystem in a safe and scaleable way?
    Minions can provide each job with its own folder for things that need to be downloaded. Cursor can monitor the size of that filesystem and report it back to the user along with the size of all of pipeline folders and the total size of the disk. User will be responsible for making sure they don't go over disk limits. This does mean that minions don't have a centralized file system that they can use and all the pros/cons that come with that.
    (how do we present this folder to the user? Does plugins have an API for this?)


-----

Pipeline object {
    id
    name
    description
    last_run
    last_run_by
    created
    modified
    depends_on
}

PipelineRun object {

}

worker object {
    id
    hostname
    state

    GetPipeline/{id}
    RunPipeline
    ChangeState
}


master object {
    GetPipeline/{can take an ID}
    ListPipelines
    GetMinion/{can take an ID}
    ListMinions
}


* Pipelines are made up of multiple jobs
* Jobs can depend on one another
* Jobs are snippets of code that do something and then return a status


CreatePipeline ---
	// attempt to compile github repo saving the resulting binary to a specific directory called plugins
	// we should only try to compile a specified folder that we determine
	// pipelines all have a unique name
