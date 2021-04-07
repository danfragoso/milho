package interpreter

import (
	"github.com/danfragoso/milho/parser"
)

func Run(nodes []*parser.Node) ([]Result, error) {
	var session *Session
	var results []Result

	var err error

	for _, node := range nodes {
		if session == nil {
			session, err = createSession(node)
		} else {
			err = updateSession(session, node)
		}

		if err != nil {
			return results, err
		}

		res, err := eval(session.Tree, session)
		if err != nil {
			return results, err
		}

		results = append(results, res)
	}

	return results, nil
}

func RunFromSession(nodes []*parser.Node, session *Session) ([]Result, error) {
	var results []Result

	for _, node := range nodes {
		err := updateSession(session, node)
		if err != nil {
			return results, err
		}

		res, err := eval(session.Tree, session)
		if err != nil {
			return results, err
		}

		results = append(results, res)
	}

	return results, nil
}
