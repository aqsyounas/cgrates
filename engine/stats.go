/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package engine

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

/*
// NewStatService initializes a StatService
func NewStatService(dataDB DataDB, ms Marshaler, storeInterval time.Duration) (ss *StatService, err error) {
	ss = &StatService{dataDB: dataDB, ms: ms, storeInterval: storeInterval,
		stopStoring: make(chan struct{})}
	sqPrfxs, err := dataDB.GetKeysForPrefix(utils.StatsConfigPrefix)
	if err != nil {
		return nil, err
	}
	go ss.dumpStoredMetrics() // start dumpStoredMetrics loop
	return
}

// StatService builds stats for events
type StatService struct {
	dataDB        DataDB
	ms            Marshaler
	storeInterval time.Duration
	stopStoring   chan struct{}
}

// ListenAndServe loops keeps the service alive
func (ss *StatService) ListenAndServe(exitChan chan bool) error {
	e := <-exitChan
	exitChan <- e // put back for the others listening for shutdown request
	return nil
}

// Called to shutdown the service
// ToDo: improve with context, ie following http implementation
func (ss *StatService) Shutdown() error {
	utils.Logger.Info("<StatS> service shutdown initialized")
	close(ss.stopStoring)
	ss.storeMetrics()
	utils.Logger.Info("<StatS> service shutdown complete")
	return nil
}

// setQueue adds or modifies a queue into cache
// sort will reorder the ss.queues
func (ss *StatService) loadQueue(qID string) (q *StatQueue, err error) {
	sq, err := ss.dataDB.GetStatsConfig(qID)
	if err != nil {
		return nil, err
	}
	return NewStatQueue(ss.evCache, ss.ms, sq, sqSM)
}

func (ss *StatService) setQueue(q *StatQueue) {
	ss.queuesCache[q.cfg.ID] = q
	ss.queues = append(ss.queues, q)
}

// remQueue will remove a queue based on it's ID
func (ss *StatService) remQueue(qID string) (si *StatQueue) {
	si = ss.queuesCache[qID]
	ss.queues.remWithID(qID)
	delete(ss.queuesCache, qID)
	return
}

// store stores the necessary storedMetrics to dataDB
func (ss *StatService) storeMetrics() {
	for _, si := range ss.queues {
		if !si.cfg.Store || !si.dirty { // no need to save
			continue
		}
		if siSM := si.GetStoredMetrics(); siSM != nil {
			if err := ss.dataDB.SetSQStoredMetrics(siSM); err != nil {
				utils.Logger.Warning(
					fmt.Sprintf("<StatService> failed saving StoredMetrics for QueueID: %s, error: %s",
						si.cfg.ID, err.Error()))
			}
		}
		// randomize the CPU load and give up thread control
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Nanosecond)
	}
	return
}

// dumpStoredMetrics regularly dumps metrics to dataDB
func (ss *StatService) dumpStoredMetrics() {
	for {
		select {
		case <-ss.stopStoring:
			return
		}
		ss.storeMetrics()
		time.Sleep(ss.storeInterval)
	}
}

// processEvent processes a StatsEvent through the queues and caches it when needed
func (ss *StatService) processEvent(ev StatsEvent) (err error) {
	evStatsID := ev.ID()
	if evStatsID == "" { // ID is mandatory
		return errors.New("missing ID field")
	}
	for _, stInst := range ss.queues {
		if err := stInst.ProcessEvent(ev); err != nil {
			utils.Logger.Warning(
				fmt.Sprintf("<StatService> QueueID: %s, ignoring event with ID: %s, error: %s",
					stInst.cfg.ID, evStatsID, err.Error()))
		}
		if stInst.cfg.Blocker {
			break
		}
	}
	return
}

// V1ProcessEvent implements StatV1 method for processing an Event
func (ss *StatService) V1ProcessEvent(ev StatsEvent, reply *string) (err error) {
	if err = ss.processEvent(ev); err == nil {
		*reply = utils.OK
	}
	return
}

// V1GetQueueIDs returns list of queue IDs configured in the service
func (ss *StatService) V1GetQueueIDs(ignored struct{}, reply *[]string) (err error) {
	if len(ss.queuesCache) == 0 {
		return utils.ErrNotFound
	}
	for k := range ss.queuesCache {
		*reply = append(*reply, k)
	}
	return
}

// V1GetStringMetrics returns the metrics as string values
func (ss *StatService) V1GetStringMetrics(queueID string, reply *map[string]string) (err error) {
	sq, has := ss.queuesCache[queueID]
	if !has {
		return utils.ErrNotFound
	}
	metrics := make(map[string]string, len(sq.sqMetrics))
	for metricID, metric := range sq.sqMetrics {
		metrics[metricID] = metric.GetStringValue("")
	}
	*reply = metrics
	return
}

// V1GetFloatMetrics returns the metrics as float64 values
func (ss *StatService) V1GetFloatMetrics(queueID string, reply *map[string]float64) (err error) {
	sq, has := ss.queuesCache[queueID]
	if !has {
		return utils.ErrNotFound
	}
	metrics := make(map[string]float64, len(sq.sqMetrics))
	for metricID, metric := range sq.sqMetrics {
		metrics[metricID] = metric.GetFloat64Value()
	}
	*reply = metrics
	return
}

// ArgsLoadQueues are the arguments passed to V1LoadQueues
type ArgsLoadQueues struct {
	QueueIDs *[]string
}

// V1LoadQueues loads the queues specified by qIDs into the service
// loads all if args.QueueIDs is nil
func (ss *StatService) V1LoadQueues(args ArgsLoadQueues, reply *string) (err error) {
	qIDs := args.QueueIDs
	if qIDs == nil {
		sqPrfxs, err := ss.dataDB.GetKeysForPrefix(utils.StatsConfigPrefix)
		if err != nil {
			return err
		}
		queueIDs := make([]string, len(sqPrfxs))
		for i, prfx := range sqPrfxs {
			queueIDs[i] = prfx[len(utils.StatsConfigPrefix):]
		}
		if len(queueIDs) != 0 {
			qIDs = &queueIDs
		}
	}
	if qIDs == nil || len(*qIDs) == 0 {
		return utils.ErrNotFound
	}
	var sQs []*StatQueue // cache here so we lock only later when data available
	for _, qID := range *qIDs {
		if _, hasPrev := ss.queuesCache[qID]; hasPrev {
			continue // don't overwrite previous, could be extended in the future by carefully checking cached events
		}
		if q, err := ss.loadQueue(qID); err != nil {
			utils.Logger.Err(fmt.Sprintf("<StatS> failed loading quueue with id: <%s>, err: <%s>",
				q.cfg.ID, err.Error()))
			continue
		} else {
			sQs = append(sQs, q)
		}
	}
	ss.Lock()
	for _, q := range sQs {
		ss.setQueue(q)
	}
	ss.queues.Sort()
	ss.Unlock()
	*reply = utils.OK
	return
}

// Call implements rpcclient.RpcClientConnection interface for internal RPC
// here for testing purposes
func (ss *StatService) Call(serviceMethod string, args interface{}, reply interface{}) error {
	methodSplit := strings.Split(serviceMethod, ".")
	if len(methodSplit) != 2 {
		return rpcclient.ErrUnsupporteServiceMethod
	}
	method := reflect.ValueOf(ss).MethodByName(methodSplit[0][len(methodSplit[0])-2:] + methodSplit[1])
	if !method.IsValid() {
		return rpcclient.ErrUnsupporteServiceMethod
	}
	params := []reflect.Value{reflect.ValueOf(args), reflect.ValueOf(reply)}
	ret := method.Call(params)
	if len(ret) != 1 {
		return utils.ErrServerError
	}
	if ret[0].Interface() == nil {
		return nil
	}
	err, ok := ret[0].Interface().(error)
	if !ok {
		return utils.ErrServerError
	}
	return err
}

*/