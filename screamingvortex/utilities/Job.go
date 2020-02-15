package utilities

import "database/sql"

type Job struct {
  Id int64 `sql:"id"`
  PercentComplete int `sql:"percent_complete"`
  Error sql.NullString `sql:"error"`
  AssetId sql.NullInt64 `sql:"asset_id"`
  step int
  totalSteps int
}

func (job *Job) TableName(rollType string) string {
  return "plan_job"
}

func (job *Job) GetId() *int64 {
  return &job.Id
}

func CreateJob(id int64, totalSteps int) *Job {
  return &Job{
    Id: id,
    PercentComplete: 0,
    Error: sql.NullString{String: "", Valid: false},
    AssetId: sql.NullInt64{Int64: 0, Valid: false},
    step: 0,
    totalSteps: totalSteps,
  }
}

func (job *Job) Step(client ClientInterface) {
  job.step++
  job.PercentComplete = 100 * job.step / job.totalSteps
  client.Update(job, "")
}

func (job *Job) Complete(client ClientInterface, assetId int64) {
  job.PercentComplete = 100
  job.AssetId.Valid = true
  job.AssetId.Int64 = assetId
  client.Update(job, "")
}

func (job *Job) Panic(client ClientInterface, msg string) {
  job.Error.Valid = true
  job.Error.String = msg
  client.Update(job, "")
}
