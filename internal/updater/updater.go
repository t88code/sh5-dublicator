package updater

//func Update(ctx context.Context,
//	sh5ClientSrc sh5api.ClientInterface,
//	sh5ClientDst sh5api.ClientInterface,
//	compareResults map[sh5api.HeadName]comparer.CompareResult,
//	headNamesForSync []sh5api.HeadName) error {
//	slog.Info("Обновление выполняется....")
//
//	for _, headName := range headNamesForSync {
//
//		switch headName {
//		case sh5api.GGroupsTree:
//
//			slog.Info(fmt.Sprint(compareResults[headName]))
//			slog.Info("compareResult.ReqRepSrc")
//			slog.Info(fmt.Sprint(compareResults[headName].ReqRepSrc))
//			slog.Info("compareResult.ReqRepDst")
//			slog.Info(fmt.Sprint(compareResults[headName].ReqRepDst))
//
//			original := []string{
//				sh5api.GUID,        // 4
//				sh5api.FIELD_209_1, // 209#1\1
//				sh5api.FIELD_3,     // 3
//				sh5api.FIELD_6,     // 6
//			}
//
//			var values [][]interface{}
//
//			l := len(compareResults[headName].ReqRepSrc.Rep.(*sh5api.GGroupsTreeRep).ShTable)
//
//			// 1 - ищем индекс в original
//			// 2 - добавляем по найденному индексу в values
//
//			// 1 делаем поиск индекс в Rep
//			// формируем []originalStruct
//			// добавить original struct = {sh5api.GUID, indexInRep}
//			// 2 формируем отдельно инсерты под каждый Value на основе []originalStruct
//
//			// ТУТ начинаем цилк перебора элементов
//			// для каждого элемента подбираем originalIndex
//			//l := len(compareResults[headName].ReqRepSrc.Rep.(*sh5api.GGroupsTreeRep).ShTable)
//			//for i, valueSrc := range compareResults[headName].ReqRepSrc.Rep.(*sh5api.GGroupsTreeRep).ShTable[l-1].Values {
//			//
//			//}
//
//			for _, originalScheme := range original {
//				var originalIndex int        // индекс original для получения value по этому индексу
//				var originalIndexIsFind bool // индекс original был найден
//
//				// поиск индекса originalIndex в rep, чтобы использовать его для получения value
//				for originalRepIndex, originalRep := range compareResults[headName].ReqRepSrc.Rep.(*sh5api.GGroupsTreeRep).ShTable[l-1].Original {
//					if originalScheme == originalRep {
//						originalIndex = originalRepIndex
//						originalIndexIsFind = true
//					}
//				}
//
//				if originalIndexIsFind {
//					// тут надо прогнать все значение на insert, modify, delete
//
//					// добавили 1й элемент
//					values = append(values, []interface{}{
//						compareResults[headName].ReqRepSrc.Rep.(*sh5api.GGroupsTreeRep).ShTable[l-1].Values[originalIndex][0],
//					})
//
//					//for _, value := range compareResults[headName].ReqRepSrc.Rep.(*sh5api.GGroupsTreeRep).ShTable[l-1].Values[originalIndex] {
//					//	// тут записать только одно значение
//					//	values = append(values, []interface{}{value})
//					//	// values = append(values, []interface{}{1}) // 4
//					//	// values = append(values, []interface{}{1}) // 209#1\1
//					//	// values = append(values, []interface{}{1}) // 3
//					//	// values = append(values, []interface{}{1}) // 6
//					//	break
//					//}
//				} else {
//					return fmt.Errorf("не найден originalIndex для поля (%s) в справочнике (%s)", originalScheme, headName)
//				}
//
//			}
//
//			// отправили в API
//			insGGroupRep, err := sh5ClientDst.InsertGGroup(ctx,
//				original,
//				values,
//			)
//			if err != nil {
//				return err
//			}
//
//			// повторяем все сначала, но уже добавляем следующий элемент
//
//			slog.Info("Обновление справочника (%s) выполнено", headName)
//			slog.Info(fmt.Sprint(insGGroupRep))
//		}
//
//	}
//
//	return nil
//}
